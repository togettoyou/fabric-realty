// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"

	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
)

// Network represents a network of nodes that are members of a specific Fabric channel. The Network can be used to
// access deployed smart contracts, and to listen for events emitted when blocks are committed to the ledger. Network
// instances are obtained from a Gateway using the [Gateway.GetNetwork] method.
//
// To safely handle connection errors during eventing, it is recommended to use a checkpointer to track eventing
// progress. This allows eventing to be resumed with no loss or duplication of events.
type Network struct {
	client             *gatewayClient
	signingID          *signingIdentity
	name               string
	tlsCertificateHash []byte
}

// Name of the Fabric channel this network represents.
func (network *Network) Name() string {
	return network.name
}

// GetContract returns a Contract representing the default smart contract for the named chaincode.
func (network *Network) GetContract(chaincodeName string) *Contract {
	return network.GetContractWithName(chaincodeName, "")
}

// GetContractWithName returns a Contract representing a smart contract within a named chaincode.
func (network *Network) GetContractWithName(chaincodeName string, contractName string) *Contract {
	return &Contract{
		client:        network.client,
		signingID:     network.signingID,
		channelName:   network.name,
		chaincodeName: chaincodeName,
		contractName:  contractName,
	}
}

// ChaincodeEventsOption implements an option for a chaincode events request.
//
// If both a start block and checkpoint are specified, and the checkpoint has a valid position set, the checkpoint
// position is used and the specified start block is ignored. If the checkpoint is unset then the start block is used.
//
// If no start position is specified, eventing begins from the next committed block.
type ChaincodeEventsOption eventOption

// ChaincodeEvents returns a channel from which chaincode events emitted by transaction functions in the specified
// chaincode can be read.
func (network *Network) ChaincodeEvents(ctx context.Context, chaincodeName string, options ...ChaincodeEventsOption) (<-chan *ChaincodeEvent, error) {
	events, err := network.NewChaincodeEventsRequest(chaincodeName, options...)
	if err != nil {
		return nil, err
	}

	return events.Events(ctx)
}

// NewChaincodeEventsRequest creates a request to read events emitted by the specified chaincode. Supports off-line
// signing flow.
func (network *Network) NewChaincodeEventsRequest(chaincodeName string, options ...ChaincodeEventsOption) (*ChaincodeEventsRequest, error) {
	builder := &chaincodeEventsBuilder{
		eventsBuilder: eventsBuilder{
			signingID:   network.signingID,
			channelName: network.name,
			client:      network.client,
		},
		chaincodeName: chaincodeName,
	}

	for _, option := range options {
		if err := option(&builder.eventsBuilder); err != nil {
			return nil, err
		}
	}

	return builder.build()
}

// BlockEventsOption implements an option for a block events request.
type BlockEventsOption eventOption

// BlockEvents returns a channel from which block events can be read.
func (network *Network) BlockEvents(ctx context.Context, options ...BlockEventsOption) (<-chan *common.Block, error) {
	events, err := network.NewBlockEventsRequest(options...)
	if err != nil {
		return nil, err
	}

	return events.Events(ctx)
}

// NewBlockEventsRequest creates a request to read block events. Supports off-line signing flow.
func (network *Network) NewBlockEventsRequest(options ...BlockEventsOption) (*BlockEventsRequest, error) {
	builder := &blockEventsBuilder{
		baseBlockEventsBuilder{
			eventsBuilder: eventsBuilder{
				signingID:   network.signingID,
				channelName: network.name,
				client:      network.client,
			},
			tlsCertificateHash: network.tlsCertificateHash,
		},
	}

	for _, option := range options {
		if err := option(&builder.eventsBuilder); err != nil {
			return nil, err
		}
	}

	return builder.build()
}

// FilteredBlockEvents returns a channel from which filtered block events can be read.
func (network *Network) FilteredBlockEvents(ctx context.Context, options ...BlockEventsOption) (<-chan *peer.FilteredBlock, error) {
	events, err := network.NewFilteredBlockEventsRequest(options...)
	if err != nil {
		return nil, err
	}

	return events.Events(ctx)
}

// NewFilteredBlockEventsRequest creates a request to read filtered block events. Supports off-line signing flow.
func (network *Network) NewFilteredBlockEventsRequest(options ...BlockEventsOption) (*FilteredBlockEventsRequest, error) {
	builder := &filteredBlockEventsBuilder{
		baseBlockEventsBuilder{
			eventsBuilder: eventsBuilder{
				signingID:   network.signingID,
				channelName: network.name,
				client:      network.client,
			},
			tlsCertificateHash: network.tlsCertificateHash,
		},
	}

	for _, option := range options {
		if err := option(&builder.eventsBuilder); err != nil {
			return nil, err
		}
	}

	return builder.build()
}

// BlockAndPrivateDataEvents returns a channel from which block and private data events can be read.
func (network *Network) BlockAndPrivateDataEvents(ctx context.Context, options ...BlockEventsOption) (<-chan *peer.BlockAndPrivateData, error) {
	events, err := network.NewBlockAndPrivateDataEventsRequest(options...)
	if err != nil {
		return nil, err
	}

	return events.Events(ctx)
}

// NewBlockAndPrivateDataEventsRequest creates a request to read block and private data events. Supports off-line signing flow.
func (network *Network) NewBlockAndPrivateDataEventsRequest(options ...BlockEventsOption) (*BlockAndPrivateDataEventsRequest, error) {
	builder := &blockAndPrivateDataEventsBuilder{
		baseBlockEventsBuilder{
			eventsBuilder: eventsBuilder{
				signingID:   network.signingID,
				channelName: network.name,
				client:      network.client,
			},
			tlsCertificateHash: network.tlsCertificateHash,
		},
	}

	for _, option := range options {
		if err := option(&builder.eventsBuilder); err != nil {
			return nil, err
		}
	}

	return builder.build()
}
