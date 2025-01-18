// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

// InMemoryCheckpointer is a non-persistent [Checkpoint] implementation. It can be used to checkpoint progress after
// successfully processing events, allowing eventing to be resumed from this point.
type InMemoryCheckpointer struct {
	blockNumber   uint64
	transactionID string
}

// CheckpointBlock records a successfully processed block.
func (c *InMemoryCheckpointer) CheckpointBlock(blockNumber uint64) {
	c.CheckpointTransaction(blockNumber+1, "")
}

// CheckpointTransaction records a successfully processed transaction within a given block.
func (c *InMemoryCheckpointer) CheckpointTransaction(blockNumber uint64, transactionID string) {
	c.blockNumber = blockNumber
	c.transactionID = transactionID
}

// CheckpointChaincodeEvent records a successfully processed chaincode event.
func (c *InMemoryCheckpointer) CheckpointChaincodeEvent(event *ChaincodeEvent) {
	c.CheckpointTransaction(event.BlockNumber, event.TransactionID)
}

// BlockNumber in which the next event is expected.
func (c *InMemoryCheckpointer) BlockNumber() uint64 {
	return c.blockNumber
}

// TransactionID of the last successfully processed event within the current block.
func (c *InMemoryCheckpointer) TransactionID() string {
	return c.transactionID
}
