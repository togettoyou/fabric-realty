// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
)

// Contract represents a smart contract, and allows applications to:
//   - Evaluate transactions that query state from the ledger using the [Contract.EvaluateTransaction] method.
//   - Submit transactions that store state to the ledger using the [Contract.SubmitTransaction] method.
//
// For more complex transaction invocations, such as including transient data, transactions can be evaluated or
// submitted using the [Contract.Evaluate] or [Contract.Submit] methods respectively. The result of a submitted
// transaction can be accessed prior to its commit to the ledger using [Contract.SubmitAsync].
//
// A finer-grained transaction flow can be employed by using [Contract.NewProposal]. This allows retry of individual
// steps in the flow in response to errors.
//
// # Off-line signing
//
// By default, proposal, transaction and commit status messages will be signed using the signing implementation
// specified when connecting the Gateway. In cases where an external client holds the signing credentials, a signing
// implementation can be omitted when connecting the Gateway and off-line signing can be carried out by:
//  1. Returning the serialized proposal, transaction or commit status message along with its digest to the client for
//     them to generate a signature.
//  2. With the serialized message and signature received from the client to create a signed proposal, transaction or
//     commit using the [Gateway.NewSignedProposal], [Gateway.NewSignedTransaction] or [Gateway.NewSignedCommit] methods
//     respectively.
//
// Note that the message digest is created with the hash implementation specified by the [WithHash] option passed to the
// [Connect] function used to create the [Gateway] instance. For off-line signing implementations that require the
// entire message content, a NONE (or no-op) hash implementation should be specified.
type Contract struct {
	client        *gatewayClient
	signingID     *signingIdentity
	channelName   string
	chaincodeName string
	contractName  string
}

// ChaincodeName of the chaincode that contains this smart contract.
func (contract *Contract) ChaincodeName() string {
	return contract.chaincodeName
}

// ContractName of the contract within the chaincode, or an empty string for the default smart contract.
func (contract *Contract) ContractName() string {
	return contract.contractName
}

// EvaluateTransaction will evaluate a transaction function and return its results. A transaction proposal will be
// evaluated on endorsing peers but the transaction will not be sent to the ordering service and so will not be
// committed to the ledger. This can be used for querying the world state.
//
// This method is equivalent to:
//
//	contract.Evaluate(name, WithArguments(args...))
func (contract *Contract) EvaluateTransaction(name string, args ...string) ([]byte, error) {
	return contract.Evaluate(name, WithArguments(args...))
}

// Evaluate a transaction function and return its result. This method provides greater control over the transaction
// proposal content and the endorsing peers on which it is evaluated. This allows transaction functions to be evaluated
// where the proposal must include transient data.
func (contract *Contract) Evaluate(transactionName string, options ...ProposalOption) ([]byte, error) {
	proposal, err := contract.NewProposal(transactionName, options...)
	if err != nil {
		return nil, err
	}

	return proposal.Evaluate()
}

// EvaluateWithContext evaluates a transaction function in the scope of a specific context and return its result. This
// method provides greater control over the transaction proposal content and the endorsing peers on which it is
// evaluated. This allows transaction functions to be evaluated where the proposal must include transient data.
func (contract *Contract) EvaluateWithContext(ctx context.Context, transactionName string, options ...ProposalOption) ([]byte, error) {
	proposal, err := contract.NewProposal(transactionName, options...)
	if err != nil {
		return nil, err
	}
	return proposal.EvaluateWithContext(ctx)
}

// SubmitTransaction will submit a transaction to the ledger and return its result only after it is committed to the
// ledger. The transaction function will be evaluated on endorsing peers and then submitted to the ordering service to
// be committed to the ledger.
//
// This method may return different error types depending on the point in the transaction invocation that a failure
// occurs. The error can be inspected with errors.Is or errors.As.
//
// This method is equivalent to:
//
//	contract.Submit(name, client.WithArguments(args...))
func (contract *Contract) SubmitTransaction(name string, args ...string) ([]byte, error) {
	return contract.Submit(name, WithArguments(args...))
}

// Submit a transaction to the ledger and return its result only after it has been committed to the ledger. This method
// provides greater control over the transaction proposal content and the endorsing peers on which it is evaluated.
// This allows transaction functions to be submitted where the proposal must include transient data.
//
// This method may return different error types depending on the point in the transaction invocation that a failure
// occurs. The error can be inspected with errors.Is or errors.As.
func (contract *Contract) Submit(transactionName string, options ...ProposalOption) ([]byte, error) {
	result, commit, err := contract.SubmitAsync(transactionName, options...)
	if err != nil {
		return result, err
	}

	status, err := commit.Status()
	if err != nil {
		return result, err
	}

	if !status.Successful {
		return nil, newCommitError(status.TransactionID, status.Code)
	}

	return result, nil
}

// SubmitWithContext submit a transaction to the ledger in the scope of a specific Context and return its result only
// after it has been committed to the ledger. This method provides greater control over the transaction proposal
// content and the endorsing peers on which it is evaluated. This allows transaction functions to be submitted where
// the proposal must include transient data.
//
// This method may return different error types depending on the point in the transaction invocation that a failure
// occurs. The error can be inspected with errors.Is or errors.As.
func (contract *Contract) SubmitWithContext(ctx context.Context, transactionName string, options ...ProposalOption) ([]byte, error) {

	result, commit, err := contract.SubmitAsyncWithContext(ctx, transactionName, options...)
	if err != nil {
		return result, err
	}

	status, err := commit.StatusWithContext(ctx)
	if err != nil {
		return result, err
	}

	if !status.Successful {
		return nil, newCommitError(status.TransactionID, status.Code)
	}

	return result, nil
}

// SubmitAsync submits a transaction to the ledger and returns its result immediately after successfully sending to the
// orderer, along with a Commit that can be used to wait for it to be committed to the ledger.
//
// This method may return different error types depending on the point in the transaction invocation that a failure
// occurs. The error can be inspected with errors.Is or errors.As.
func (contract *Contract) SubmitAsync(transactionName string, options ...ProposalOption) ([]byte, *Commit, error) {
	proposal, err := contract.NewProposal(transactionName, options...)
	if err != nil {
		return nil, nil, err
	}

	transaction, err := proposal.Endorse()
	if err != nil {
		return nil, nil, err
	}

	result := transaction.Result()

	commit, err := transaction.Submit()
	if err != nil {
		return result, nil, err
	}

	return result, commit, nil
}

// SubmitAsyncWithContext submits a transaction to the ledger in the scope of a specific context and returns its result
// immediately after successfully sending to the orderer, along with a Commit that can be used to wait for it to be
// committed to the ledger.
//
// This method may return different error types depending on the point in the transaction invocation that a failure
// occurs. The error can be inspected with errors.Is or errors.As.
func (contract *Contract) SubmitAsyncWithContext(ctx context.Context, transactionName string, options ...ProposalOption) ([]byte, *Commit, error) {
	proposal, err := contract.NewProposal(transactionName, options...)
	if err != nil {
		return nil, nil, err
	}

	transaction, err := proposal.EndorseWithContext(ctx)
	if err != nil {
		return nil, nil, err
	}

	result := transaction.Result()

	commit, err := transaction.SubmitWithContext(ctx)
	if err != nil {
		return result, nil, err
	}

	return result, commit, nil
}

// NewProposal creates a proposal that can be sent to peers for endorsement. Supports off-line signing transaction flow.
func (contract *Contract) NewProposal(transactionName string, options ...ProposalOption) (*Proposal, error) {
	builder, err := newProposalBuilder(
		contract.client,
		contract.signingID,
		contract.channelName,
		contract.chaincodeName,
		contract.qualifiedTransactionName(transactionName),
	)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		if err := option(builder); err != nil {
			return nil, err
		}
	}

	return builder.build()
}

func (contract *Contract) qualifiedTransactionName(name string) string {
	if len(contract.contractName) > 0 {
		return contract.contractName + ":" + name
	}
	return name
}
