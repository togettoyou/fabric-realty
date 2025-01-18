// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package identity

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"fmt"
)

// Sign function generates a digital signature of the supplied digest.
type Sign = func(digest []byte) ([]byte, error)

// NewPrivateKeySign returns a Sign function that uses the supplied private key.
//
// Currently supported private key types are:
//   - ECDSA.
//   - Ed25519.
//
// Note that the Sign implementations have different expectations on the input data supplied to them.
//
// The ECDSA signers operate on a pre-computed message digest, and should be combined with an appropriate hash
// algorithm. P-256 is typically used with a SHA-256 hash, and P-384 is typically used with a SHA-384 hash.
//
// The Ed25519 signer operates on the full message content, and should be combined with a NONE (or no-op) hash
// implementation to ensure the complete message is passed to the signer.
func NewPrivateKeySign(privateKey crypto.PrivateKey) (Sign, error) {
	switch key := privateKey.(type) {
	case *ecdsa.PrivateKey:
		return ecdsaPrivateKeySign(key), nil
	case ed25519.PrivateKey:
		return ed25519PrivateKeySign(key), nil
	default:
		return nil, fmt.Errorf("unsupported key type: %T", privateKey)
	}
}

func ed25519PrivateKeySign(privateKey ed25519.PrivateKey) Sign {
	return func(message []byte) ([]byte, error) {
		signature := ed25519.Sign(privateKey, message)
		return signature, nil
	}
}
