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

// Proposal represents a transaction proposal that can be sent to peers for endorsement or evaluated as a query.
type Proposal struct {
	client              *gatewayClient
	signingID           *signingIdentity
	channelID           string
	proposedTransaction *gateway.ProposedTransaction
}

// Bytes of the serialized proposal message.
func (proposal *Proposal) Bytes() ([]byte, error) {
	transactionBytes, err := proto.Marshal(proposal.proposedTransaction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall Proposal protobuf: %w", err)
	}

	return transactionBytes, nil
}

// Digest of the proposal. This is used to generate a digital signature.
func (proposal *Proposal) Digest() []byte {
	return proposal.signingID.Hash(proposal.proposedTransaction.Proposal.ProposalBytes)
}

// TransactionID for the proposal.
func (proposal *Proposal) TransactionID() string {
	return proposal.proposedTransaction.GetTransactionId()
}

// Endorse the proposal and obtain an endorsed transaction for submission to the orderer.
func (proposal *Proposal) Endorse(opts ...grpc.CallOption) (*Transaction, error) {
	return proposal.endorse(proposal.client.Endorse, opts...)
}

// EndorseWithContext uses ths supplied context to endorse the proposal and obtain an endorsed transaction for
// submission to the orderer.
func (proposal *Proposal) EndorseWithContext(ctx context.Context, opts ...grpc.CallOption) (*Transaction, error) {
	return proposal.endorse(
		func(in *gateway.EndorseRequest, opts ...grpc.CallOption) (*gateway.EndorseResponse, error) {
			return proposal.client.EndorseWithContext(ctx, in, opts...)
		},
		opts...,
	)
}

func (proposal *Proposal) endorse(
	call func(in *gateway.EndorseRequest, opts ...grpc.CallOption) (*gateway.EndorseResponse, error),
	opts ...grpc.CallOption,
) (*Transaction, error) {
	if err := proposal.sign(); err != nil {
		return nil, err
	}

	endorseRequest := &gateway.EndorseRequest{
		TransactionId:          proposal.proposedTransaction.GetTransactionId(),
		ChannelId:              proposal.channelID,
		ProposedTransaction:    proposal.proposedTransaction.GetProposal(),
		EndorsingOrganizations: proposal.proposedTransaction.GetEndorsingOrganizations(),
	}
	response, err := call(endorseRequest, opts...)
	if err != nil {
		return nil, err
	}

	preparedTransaction := &gateway.PreparedTransaction{
		TransactionId: proposal.proposedTransaction.GetTransactionId(),
		Envelope:      response.GetPreparedTransaction(),
	}
	return newTransaction(proposal.client, proposal.signingID, preparedTransaction)
}

// Evaluate the proposal and obtain a transaction result. This is effectively a query.
func (proposal *Proposal) Evaluate(opts ...grpc.CallOption) ([]byte, error) {
	return proposal.evaluate(proposal.client.Evaluate, opts...)
}

// EvaluateWithContext uses ths supplied context to evaluate the proposal and obtain a transaction result. This is
// effectively a query.
func (proposal *Proposal) EvaluateWithContext(ctx context.Context, opts ...grpc.CallOption) ([]byte, error) {
	return proposal.evaluate(
		func(in *gateway.EvaluateRequest, opts ...grpc.CallOption) (*gateway.EvaluateResponse, error) {
			return proposal.client.EvaluateWithContext(ctx, in, opts...)
		},
		opts...,
	)
}

func (proposal *Proposal) evaluate(
	call func(in *gateway.EvaluateRequest, opts ...grpc.CallOption) (*gateway.EvaluateResponse, error),
	opts ...grpc.CallOption,
) ([]byte, error) {
	if err := proposal.sign(); err != nil {
		return nil, err
	}

	evaluateRequest := &gateway.EvaluateRequest{
		TransactionId:       proposal.proposedTransaction.GetTransactionId(),
		ChannelId:           proposal.channelID,
		ProposedTransaction: proposal.proposedTransaction.GetProposal(),
		TargetOrganizations: proposal.proposedTransaction.GetEndorsingOrganizations(),
	}
	response, err := call(evaluateRequest, opts...)
	if err != nil {
		return nil, err
	}

	return response.GetResult().GetPayload(), nil
}

func (proposal *Proposal) setSignature(signature []byte) {
	proposal.proposedTransaction.Proposal.Signature = signature
}

func (proposal *Proposal) isSigned() bool {
	return len(proposal.proposedTransaction.GetProposal().GetSignature()) > 0
}

func (proposal *Proposal) sign() error {
	if proposal.isSigned() {
		return nil
	}

	digest := proposal.Digest()
	signature, err := proposal.signingID.Sign(digest)
	if err != nil {
		return err
	}

	proposal.setSignature(signature)

	return nil
}
