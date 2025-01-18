// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"encoding/json"
	"os"
)

// FileCheckpointer is a [Checkpoint] implementation backed by persistent file storage. It can be used to checkpoint
// progress after successfully processing events, allowing eventing to be resumed from this point.
//
// Instances should be created using the [NewFileCheckpointer] constructor function. [FileCheckpointer.Close] should
// be called when the checkpointer is no longer needed to free resources.
type FileCheckpointer struct {
	file  *os.File
	state *checkpointState
}

type checkpointState struct {
	BlockNumber   uint64 `json:"blockNumber"`
	TransactionID string `json:"transactionId"`
}

// NewFileCheckpointer creates a properly initialized FileCheckpointer.
func NewFileCheckpointer(name string) (*FileCheckpointer, error) {
	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0600) //#nosec G304 -- Caller responsible for safe file name
	if err != nil {
		return nil, err
	}

	checkpointer := &FileCheckpointer{
		file:  file,
		state: &checkpointState{},
	}

	if fileInfo, err := file.Stat(); err == nil && fileInfo.Size() > 0 {
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(checkpointer.state); err != nil {
			return nil, err
		}
	}

	if err := checkpointer.save(); err != nil {
		return nil, err
	}

	return checkpointer, nil
}

// CheckpointBlock records a successfully processed block.
func (c *FileCheckpointer) CheckpointBlock(blockNumber uint64) error {
	return c.CheckpointTransaction(blockNumber+1, "")
}

// CheckpointTransaction records a successfully processed transaction within a given block.
func (c *FileCheckpointer) CheckpointTransaction(blockNumber uint64, transactionID string) error {
	c.state.BlockNumber = blockNumber
	c.state.TransactionID = transactionID
	return c.save()
}

// CheckpointChaincodeEvent records a successfully processed chaincode event.
func (c *FileCheckpointer) CheckpointChaincodeEvent(event *ChaincodeEvent) error {
	return c.CheckpointTransaction(event.BlockNumber, event.TransactionID)
}

// BlockNumber in which the next event is expected.
func (c *FileCheckpointer) BlockNumber() uint64 {
	return c.state.BlockNumber
}

// TransactionID of the last successfully processed event within the current block.
func (c *FileCheckpointer) TransactionID() string {
	return c.state.TransactionID
}

// Close the checkpointer when it is no longer needed to free resources.
func (c *FileCheckpointer) Close() error {
	return c.file.Close()
}

// Sync commits the current state to stable storage.
func (c *FileCheckpointer) Sync() error {
	return c.file.Sync()
}

func (c *FileCheckpointer) save() error {
	data, err := json.Marshal(c.state)
	if err != nil {
		return err
	}

	size, err := c.file.WriteAt(data, 0)
	if err != nil {
		return err
	}

	return c.file.Truncate(int64(size))
}
