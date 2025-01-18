// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func newTransaction(client *gatewayClient, signingID *signingIdentity, preparedTransaction *gateway.PreparedTransaction) (*Transaction, error) {
	txInfo, err := parseTransactionEnvelope(preparedTransaction.GetEnvelope())
	if err != nil {
		return nil, err
	}

	transaction := &Transaction{
		client:              client,
		signingID:           signingID,
		channelID:           txInfo.ChannelName,
		preparedTransaction: preparedTransaction,
		result:              txInfo.Result,
	}
	return transaction, nil
}

// Transaction represents an endorsed transaction that can be submitted to the orderer for commit to the ledger.
type Transaction struct {
	client              *gatewayClient
	signingID           *signingIdentity
	channelID           string
	preparedTransaction *gateway.PreparedTransaction
	result              []byte
}

// Result of the proposed transaction invocation.
func (transaction *Transaction) Result() []byte {
	return transaction.result
}

// Bytes of the serialized transaction.
func (transaction *Transaction) Bytes() ([]byte, error) {
	transactionBytes, err := proto.Marshal(transaction.preparedTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall PreparedTransaction protobuf: %w", err)
	}

	return transactionBytes, nil
}

// Digest of the transaction. This is used to generate a digital signature.
func (transaction *Transaction) Digest() []byte {
	return transaction.signingID.Hash(transaction.preparedTransaction.GetEnvelope().GetPayload())
}

// TransactionID of the transaction.
func (transaction *Transaction) TransactionID() string {
	return transaction.preparedTransaction.GetTransactionId()
}

// Submit the transaction to the orderer for commit to the ledger.
func (transaction *Transaction) Submit(opts ...grpc.CallOption) (*Commit, error) {
	return transaction.submit(transaction.client.Submit, opts...)
}

// SubmitWithContext uses the supplied context to submit the transaction to the orderer for commit to the ledger.
func (transaction *Transaction) SubmitWithContext(ctx context.Context, opts ...grpc.CallOption) (*Commit, error) {
	return transaction.submit(
		func(in *gateway.SubmitRequest, opts ...grpc.CallOption) (*gateway.SubmitResponse, error) {
			return transaction.client.SubmitWithContext(ctx, in, opts...)
		},
		opts...,
	)
}

func (transaction *Transaction) submit(
	call func(in *gateway.SubmitRequest, opts ...grpc.CallOption) (*gateway.SubmitResponse, error),
	opts ...grpc.CallOption,
) (*Commit, error) {
	if err := transaction.sign(); err != nil {
		return nil, err
	}

	// Build before the submit to avoid chance of errors after the submit
	statusRequest, err := transaction.newSignedCommitStatusRequest()
	if err != nil {
		return nil, err
	}

	submitRequest := &gateway.SubmitRequest{
		TransactionId:       transaction.TransactionID(),
		ChannelId:           transaction.channelID,
		PreparedTransaction: transaction.preparedTransaction.GetEnvelope(),
	}
	_, err = call(submitRequest, opts...)
	if err != nil {
		return nil, err
	}

	return newCommit(transaction.client, transaction.signingID, transaction.TransactionID(), statusRequest), nil
}

func (transaction *Transaction) sign() error {
	if transaction.isSigned() {
		return nil
	}

	digest := transaction.Digest()
	signature, err := transaction.signingID.Sign(digest)
	if err != nil {
		return err
	}

	transaction.setSignature(signature)

	return nil
}

func (transaction *Transaction) isSigned() bool {
	return len(transaction.preparedTransaction.GetEnvelope().GetSignature()) > 0
}

func (transaction *Transaction) setSignature(signature []byte) {
	transaction.preparedTransaction.Envelope.Signature = signature
}

func (transaction *Transaction) newSignedCommitStatusRequest() (*gateway.SignedCommitStatusRequest, error) {
	creator, err := transaction.signingID.Creator()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize identity: %w", err)
	}

	request := &gateway.CommitStatusRequest{
		ChannelId:     transaction.channelID,
		TransactionId: transaction.TransactionID(),
		Identity:      creator,
	}

	requestBytes, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	signedRequest := &gateway.SignedCommitStatusRequest{
		Request: requestBytes,
	}
	return signedRequest, nil
}
