// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type proposalBuilder struct {
	client          *gatewayClient
	signingID       *signingIdentity
	channelName     string
	chaincodeName   string
	transactionName string
	transactionCtx  *transactionContext
	transient       map[string][]byte
	endorsingOrgs   []string
	args            [][]byte
}

func newProposalBuilder(
	client *gatewayClient,
	signingID *signingIdentity,
	channelName string,
	chaincodeName string,
	transactionName string,
) (*proposalBuilder, error) {
	transactionCtx, err := newTransactionContext(signingID)
	if err != nil {
		return nil, err
	}

	builder := &proposalBuilder{
		client:          client,
		signingID:       signingID,
		channelName:     channelName,
		chaincodeName:   chaincodeName,
		transactionName: transactionName,
		transactionCtx:  transactionCtx,
	}
	return builder, nil
}

func (builder *proposalBuilder) build() (*Proposal, error) {
	proposalBytes, err := builder.proposalBytes()
	if err != nil {
		return nil, err
	}

	proposal := &Proposal{
		client:    builder.client,
		signingID: builder.signingID,
		channelID: builder.channelName,
		proposedTransaction: &gateway.ProposedTransaction{
			TransactionId: builder.transactionCtx.TransactionID,
			Proposal: &peer.SignedProposal{
				ProposalBytes: proposalBytes,
			},
			EndorsingOrganizations: builder.endorsingOrgs,
		},
	}
	return proposal, nil
}

func (builder *proposalBuilder) proposalBytes() ([]byte, error) {
	headerBytes, err := builder.headerBytes()
	if err != nil {
		return nil, err
	}

	chaincodeProposalBytes, err := builder.chaincodeProposalPayloadBytes()
	if err != nil {
		return nil, err
	}

	proposal := &peer.Proposal{
		Header:  headerBytes,
		Payload: chaincodeProposalBytes,
	}
	return proto.Marshal(proposal)
}

func (builder *proposalBuilder) headerBytes() ([]byte, error) {
	channelHeaderBytes, err := builder.channelHeaderBytes()
	if err != nil {
		return nil, err
	}

	signatureHeaderBytes, err := proto.Marshal(builder.transactionCtx.SignatureHeader)
	if err != nil {
		return nil, err
	}

	header := &common.Header{
		ChannelHeader:   channelHeaderBytes,
		SignatureHeader: signatureHeaderBytes,
	}
	return proto.Marshal(header)
}

func (builder *proposalBuilder) channelHeaderBytes() ([]byte, error) {
	extensionBytes, err := proto.Marshal(&peer.ChaincodeHeaderExtension{
		ChaincodeId: &peer.ChaincodeID{
			Name: builder.chaincodeName,
		},
	})
	if err != nil {
		return nil, err
	}

	channelHeader := &common.ChannelHeader{
		Type:      int32(common.HeaderType_ENDORSER_TRANSACTION),
		Timestamp: timestamppb.Now(),
		ChannelId: builder.channelName,
		TxId:      builder.transactionCtx.TransactionID,
		Epoch:     0,
		Extension: extensionBytes,
	}
	return proto.Marshal(channelHeader)
}

func (builder *proposalBuilder) chaincodeProposalPayloadBytes() ([]byte, error) {
	invocationSpecBytes, err := proto.Marshal(&peer.ChaincodeInvocationSpec{
		ChaincodeSpec: &peer.ChaincodeSpec{
			ChaincodeId: &peer.ChaincodeID{
				Name: builder.chaincodeName,
			},
			Input: &peer.ChaincodeInput{
				Args: builder.chaincodeArgs(),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	chaincodeProposalPayload := &peer.ChaincodeProposalPayload{
		Input:        invocationSpecBytes,
		TransientMap: builder.transient,
	}
	return proto.Marshal(chaincodeProposalPayload)
}

func (builder *proposalBuilder) chaincodeArgs() [][]byte {
	result := make([][]byte, len(builder.args)+1)

	result[0] = []byte(builder.transactionName)
	copy(result[1:], builder.args)

	return result
}

// ProposalOption implements an option for a transaction proposal.
type ProposalOption = func(builder *proposalBuilder) error

// WithBytesArguments appends to the transaction function arguments associated with a transaction proposal.
func WithBytesArguments(args ...[]byte) ProposalOption {
	return func(builder *proposalBuilder) error {
		builder.args = append(builder.args, args...)
		return nil
	}
}

// WithArguments appends to the transaction function arguments associated with a transaction proposal.
func WithArguments(args ...string) ProposalOption {
	return WithBytesArguments(stringsAsBytes(args)...)
}

func stringsAsBytes(strings []string) [][]byte {
	results := make([][]byte, 0, len(strings))

	for _, v := range strings {
		results = append(results, []byte(v))
	}

	return results
}

// WithTransient specifies the transient data associated with a transaction proposal.
// This is usually used in combination with [WithEndorsingOrganizations] for private data scenarios.
func WithTransient(transient map[string][]byte) ProposalOption {
	return func(builder *proposalBuilder) error {
		builder.transient = transient
		return nil
	}
}

// WithEndorsingOrganizations specifies the organizations that should endorse the transaction proposal.
// No other organizations will be sent the proposal.  This is usually used in combination with [WithTransient]
// for private data scenarios, or for state-based endorsement when specific organizations have to endorse the proposal.
func WithEndorsingOrganizations(mspids ...string) ProposalOption {
	return func(builder *proposalBuilder) error {
		builder.endorsingOrgs = mspids
		return nil
	}
}
