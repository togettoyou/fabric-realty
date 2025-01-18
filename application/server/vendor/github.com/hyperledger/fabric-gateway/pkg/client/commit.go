// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"

	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// Commit provides access to a committed transaction.
type Commit struct {
	client        *gatewayClient
	signingID     *signingIdentity
	transactionID string
	signedRequest *gateway.SignedCommitStatusRequest
}

func newCommit(
	client *gatewayClient,
	signingID *signingIdentity,
	transactionID string,
	signedRequest *gateway.SignedCommitStatusRequest,
) *Commit {
	return &Commit{
		client:        client,
		signingID:     signingID,
		transactionID: transactionID,
		signedRequest: signedRequest,
	}
}

// Bytes of the serialized commit.
func (commit *Commit) Bytes() ([]byte, error) {
	requestBytes, err := proto.Marshal(commit.signedRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall SignedCommitStatusRequest protobuf: %w", err)
	}

	return requestBytes, nil
}

// Digest of the commit status request. This is used to generate a digital signature.
func (commit *Commit) Digest() []byte {
	return commit.signingID.Hash(commit.signedRequest.GetRequest())
}

// TransactionID of the transaction.
func (commit *Commit) TransactionID() string {
	return commit.transactionID
}

// Status of the committed transaction. If the transaction has not yet committed, this call blocks until the commit
// occurs.
func (commit *Commit) Status(opts ...grpc.CallOption) (*Status, error) {
	return commit.status(commit.client.CommitStatus, opts...)
}

// StatusWithContext uses the supplied context to get the status of the committed transaction. If the transaction has
// not yet committed, this call blocks until the commit occurs.
func (commit *Commit) StatusWithContext(ctx context.Context, opts ...grpc.CallOption) (*Status, error) {
	return commit.status(
		func(in *gateway.SignedCommitStatusRequest, opts ...grpc.CallOption) (*gateway.CommitStatusResponse, error) {
			return commit.client.CommitStatusWithContext(ctx, in, opts...)
		},
		opts...,
	)
}

func (commit *Commit) status(
	call func(in *gateway.SignedCommitStatusRequest, opts ...grpc.CallOption) (*gateway.CommitStatusResponse, error),
	opts ...grpc.CallOption,
) (*Status, error) {
	if err := commit.sign(); err != nil {
		return nil, err
	}

	response, err := call(commit.signedRequest, opts...)
	if err != nil {
		return nil, err
	}

	status := &Status{
		Code:          response.GetResult(),
		Successful:    response.GetResult() == peer.TxValidationCode_VALID,
		TransactionID: commit.transactionID,
		BlockNumber:   response.GetBlockNumber(),
	}
	return status, nil
}

func (commit *Commit) sign() error {
	if commit.isSigned() {
		return nil
	}

	digest := commit.Digest()
	signature, err := commit.signingID.Sign(digest)
	if err != nil {
		return err
	}

	commit.setSignature(signature)

	return nil
}

func (commit *Commit) isSigned() bool {
	return len(commit.signedRequest.GetSignature()) > 0
}

func (commit *Commit) setSignature(signature []byte) {
	commit.signedRequest.Signature = signature
}

// Status of a committed transaction.
type Status struct {
	Code          peer.TxValidationCode
	Successful    bool
	TransactionID string
	BlockNumber   uint64
}
