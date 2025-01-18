// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build pkcs11
// +build pkcs11

package identity

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/miekg/pkcs11"
)

// HSMSignerOptions are the options required for HSM Login.
type HSMSignerOptions struct {
	Label      string
	Pin        string
	Identifier string
	UserType   int
}

// HSMSignerFactory is used to create HSM signers.
type HSMSignerFactory struct {
	ctx *pkcs11.Ctx
}

// HSMSignClose closes an HSM signer when it is no longer needed.
type HSMSignClose = func() error

// NewHSMSignerFactory creates a new HSMSignerFactory. A single factory instance should be used to create all HSM
// signers.
func NewHSMSignerFactory(library string) (*HSMSignerFactory, error) {
	if library == "" {
		return nil, fmt.Errorf("library path not provided")
	}

	ctx := pkcs11.New(library)
	if ctx == nil {
		return nil, fmt.Errorf("instantiation failed for %s", library)
	}

	if err := ctx.Initialize(); err != nil {
		return nil, fmt.Errorf("initialize failed: %w", err)
	}

	return &HSMSignerFactory{ctx}, nil
}

// NewHSMSigner creates a new HSM signer, and a close function that should be invoked when the signer is no longer
// needed. The signer implementation is thread safe but HSM operations are synchronized for a given signer so have the
// potential to become a bottleneck under load. For high volume applications, it might be beneficial to use a pool of
// HSM signers.
func (factory *HSMSignerFactory) NewHSMSigner(options HSMSignerOptions) (Sign, HSMSignClose, error) {
	if options.Label == "" {
		return nil, nil, fmt.Errorf("no Label provided")
	}

	if options.Pin == "" {
		return nil, nil, fmt.Errorf("no Pin provided")
	}

	if options.Identifier == "" {
		return nil, nil, fmt.Errorf("no Identifier provided")
	}

	slot, err := factory.findSlotForLabel(options.Label)
	if err != nil {
		return nil, nil, err
	}

	session, err := factory.createSession(slot, options.Pin)
	if err != nil {
		return nil, nil, err
	}

	privateKeyHandle, err := factory.findObjectInHSM(session, pkcs11.CKO_PRIVATE_KEY, options.Identifier)
	if err != nil {
		_ = factory.ctx.CloseSession(session)
		return nil, nil, err
	}

	signer := &hsmSigner{
		ctx:              factory.ctx,
		session:          session,
		privateKeyHandle: privateKeyHandle,
	}
	return signer.Sign, signer.Close, nil
}

// Dispose of resources held by the factory when it is no longer needed.
func (factory *HSMSignerFactory) Dispose() {
	_ = factory.ctx.Finalize()
}

func (factory *HSMSignerFactory) findSlotForLabel(label string) (uint, error) {
	slots, err := factory.ctx.GetSlotList(true)
	if err != nil {
		return 0, fmt.Errorf("get slot list failed: %w", err)
	}

	for _, slot := range slots {
		tokenInfo, err := factory.ctx.GetTokenInfo(slot)
		if err == nil && label == tokenInfo.Label {
			return slot, nil
		}
	}

	return 0, fmt.Errorf("could not find token with label %s", label)
}

func (factory *HSMSignerFactory) findObjectInHSM(session pkcs11.SessionHandle, keyType uint, identifier string) (pkcs11.ObjectHandle, error) {
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, keyType),
		pkcs11.NewAttribute(pkcs11.CKA_ID, identifier),
	}
	if err := factory.ctx.FindObjectsInit(session, template); err != nil {
		return 0, fmt.Errorf("findObjectsInit failed: %w", err)
	}
	defer func() {
		_ = factory.ctx.FindObjectsFinal(session)
	}()

	// single session instance, assume one hit only
	objs, _, err := factory.ctx.FindObjects(session, 1)
	if err != nil {
		return 0, fmt.Errorf("findObjects failed: %w", err)
	}

	if len(objs) == 0 {
		return 0, fmt.Errorf("HSM Object not found for key [%s]", hex.EncodeToString([]byte(identifier)))
	}

	return objs[0], nil
}

func (factory *HSMSignerFactory) createSession(slot uint, pin string) (pkcs11.SessionHandle, error) {
	session, err := factory.ctx.OpenSession(slot, pkcs11.CKF_SERIAL_SESSION)
	if err != nil {
		return 0, fmt.Errorf("open session failed: %w", err)
	}

	if err := factory.ctx.Login(session, pkcs11.CKU_USER, pin); err != nil && err != pkcs11.Error(pkcs11.CKR_USER_ALREADY_LOGGED_IN) {
		_ = factory.ctx.CloseSession(session)
		return 0, fmt.Errorf("login failed: %w", err)
	}

	return session, nil
}

type hsmSigner struct {
	ctx              *pkcs11.Ctx
	lock             sync.Mutex
	session          pkcs11.SessionHandle
	privateKeyHandle pkcs11.ObjectHandle
}

func (signer *hsmSigner) Close() error {
	signer.lock.Lock()
	defer signer.lock.Unlock()

	return signer.ctx.CloseSession(signer.session)
}

func (signer *hsmSigner) Sign(digest []byte) ([]byte, error) {
	signature, err := signer.hsmSign(digest)
	if err != nil {
		return nil, err
	}

	r, s := unmarshalConcatSignature(signature)

	// Only Elliptic of 256 byte keys are supported
	s = canonicalECDSASignatureSValue(s, elliptic.P256().Params().N)

	return asn1ECDSASignature(r, s)
}

func unmarshalConcatSignature(signature []byte) (r *big.Int, s *big.Int) {
	sIndex := len(signature) / 2
	r = new(big.Int).SetBytes(signature[0:sIndex])
	s = new(big.Int).SetBytes(signature[sIndex:])
	return
}

func (signer *hsmSigner) hsmSign(digest []byte) ([]byte, error) {
	signer.lock.Lock()
	defer signer.lock.Unlock()

	if err := signer.ctx.SignInit(
		signer.session,
		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)},
		signer.privateKeyHandle,
	); err != nil {
		return nil, fmt.Errorf("sign initialize failed: %w", err)
	}

	signature, err := signer.ctx.Sign(signer.session, digest)
	if err != nil {
		return nil, fmt.Errorf("sign failed: %w", err)
	}

	return signature, nil
}
