// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"

	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"google.golang.org/grpc/status"
)

type grpcError struct {
	error
}

func (e *grpcError) GRPCStatus() *status.Status {
	return status.Convert(e.error)
}

func (e *grpcError) Unwrap() error {
	return e.error
}

func newTransactionError(err error, transactionID string) *TransactionError {
	if err == nil {
		return nil
	}

	return &TransactionError{
		grpcError:     &grpcError{err},
		TransactionID: transactionID,
	}
}

// TransactionError represents an error invoking a transaction. This is a gRPC [status] error.
type TransactionError struct {
	*grpcError
	TransactionID string
}

// EndorseError represents a failure endorsing a transaction proposal.
type EndorseError struct {
	*TransactionError
}

// SubmitError represents a failure submitting an endorsed transaction to the orderer.
type SubmitError struct {
	*TransactionError
}

// CommitStatusError represents a failure obtaining the commit status of a transaction.
type CommitStatusError struct {
	*TransactionError
}

func newCommitError(transactionID string, code peer.TxValidationCode) error {
	return &CommitError{
		message:       fmt.Sprintf("transaction %s failed to commit with status code %d (%s)", transactionID, int32(code), peer.TxValidationCode_name[int32(code)]),
		TransactionID: transactionID,
		Code:          code,
	}
}

// CommitError represents a transaction that fails to commit successfully.
type CommitError struct {
	message       string
	TransactionID string
	Code          peer.TxValidationCode
}

func (e *CommitError) Error() string {
	return e.message
}
