// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package identity

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"math/big"
)

func ecdsaPrivateKeySign(privateKey *ecdsa.PrivateKey) Sign {
	n := privateKey.Params().Params().N

	return func(digest []byte) ([]byte, error) {
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, digest)
		if err != nil {
			return nil, err
		}

		s = canonicalECDSASignatureSValue(s, n)

		return asn1ECDSASignature(r, s)
	}
}

func canonicalECDSASignatureSValue(s *big.Int, curveN *big.Int) *big.Int {
	halfOrder := new(big.Int).Rsh(curveN, 1)
	if s.Cmp(halfOrder) <= 0 {
		return s
	}

	// Set s to N - s so it is in the lower part of signature space, less or equal to half order
	return new(big.Int).Sub(curveN, s)
}

type ecdsaSignature struct {
	R, S *big.Int
}

func asn1ECDSASignature(r, s *big.Int) ([]byte, error) {
	return asn1.Marshal(ecdsaSignature{
		R: r,
		S: s,
	})
}
