// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"github.com/hyperledger/fabric-protos-go-apiv2/orderer"
)

type eventsBuilder struct {
	client             *gatewayClient
	signingID          *signingIdentity
	channelName        string
	startPosition      *orderer.SeekPosition
	afterTransactionID string
}

func (builder *eventsBuilder) getStartPosition() *orderer.SeekPosition {
	if builder.startPosition != nil {
		return builder.startPosition
	}

	return &orderer.SeekPosition{
		Type: &orderer.SeekPosition_NextCommit{
			NextCommit: &orderer.SeekNextCommit{},
		},
	}
}

type eventOption = func(builder *eventsBuilder) error

// Checkpoint provides the current position for event processing.
type Checkpoint interface {
	// BlockNumber in which the next event is expected.
	BlockNumber() uint64
	// TransactionID of the last successfully processed event within the current block.
	TransactionID() string
}

// WithStartBlock reads events starting at the specified block number.
func WithStartBlock(blockNumber uint64) eventOption {
	return func(builder *eventsBuilder) error {
		builder.startPosition = &orderer.SeekPosition{
			Type: &orderer.SeekPosition_Specified{
				Specified: &orderer.SeekSpecified{
					Number: blockNumber,
				},
			},
		}
		return nil
	}
}

// WithCheckpoint reads events starting at the checkpoint position. This can be used to resume a previous eventing
// session. The zero value is ignored and a start position specified by other options or the default position is used.
func WithCheckpoint(checkpoint Checkpoint) eventOption {
	return func(builder *eventsBuilder) error {
		blockNumber := checkpoint.BlockNumber()
		transactionID := checkpoint.TransactionID()

		if blockNumber == 0 && transactionID == "" {
			return nil
		}

		builder.startPosition = &orderer.SeekPosition{
			Type: &orderer.SeekPosition_Specified{
				Specified: &orderer.SeekSpecified{
					Number: blockNumber,
				},
			},
		}
		builder.afterTransactionID = transactionID

		return nil
	}
}
