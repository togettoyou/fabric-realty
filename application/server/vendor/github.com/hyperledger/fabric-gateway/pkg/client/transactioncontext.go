// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/hyperledger/fabric-gateway/pkg/hash"
	"github.com/hyperledger/fabric-protos-go-apiv2/common"
)

type transactionContext struct {
	TransactionID   string
	SignatureHeader *common.SignatureHeader
}

func newTransactionContext(signingIdentity *signingIdentity) (*transactionContext, error) {
	nonce := make([]byte, 24)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	creator, err := signingIdentity.Creator()
	if err != nil {
		return nil, err
	}

	saltedCreator := append(nonce, creator...)
	rawTransactionID := hash.SHA256(saltedCreator)
	transactionID := hex.EncodeToString(rawTransactionID)

	signatureHeader := &common.SignatureHeader{
		Creator: creator,
		Nonce:   nonce,
	}

	transactionCtx := &transactionContext{
		TransactionID:   transactionID,
		SignatureHeader: signatureHeader,
	}
	return transactionCtx, nil
}
