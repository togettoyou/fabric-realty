// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"

	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"google.golang.org/protobuf/proto"
)

type transactionInfo struct {
	ChannelName string
	Result      []byte
}

func parseTransactionEnvelope(envelope *common.Envelope) (*transactionInfo, error) {
	payload := &common.Payload{}
	if err := proto.Unmarshal(envelope.GetPayload(), payload); err != nil {
		return nil, fmt.Errorf("failed to deserialize payload: %w", err)
	}

	channelName, err := parseChannelNameFromHeader(payload.GetHeader())
	if err != nil {
		return nil, err
	}

	result, err := parseResultFromPayload(payload)
	if err != nil {
		return nil, err
	}

	txInfo := &transactionInfo{
		ChannelName: channelName,
		Result:      result,
	}
	return txInfo, nil
}

func parseChannelNameFromHeader(header *common.Header) (string, error) {
	channelHeader := &common.ChannelHeader{}
	if err := proto.Unmarshal(header.GetChannelHeader(), channelHeader); err != nil {
		return "", fmt.Errorf("failed to deserialize channel header: %w", err)
	}

	return channelHeader.GetChannelId(), nil
}

func parseResultFromPayload(payload *common.Payload) ([]byte, error) {
	transaction := &peer.Transaction{}
	if err := proto.Unmarshal(payload.GetData(), transaction); err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %w", err)
	}

	errors := make([]error, 0)

	for _, transactionAction := range transaction.GetActions() {
		result, err := parseResultFromTransactionAction(transactionAction)
		if err == nil {
			return result, nil
		}

		errors = append(errors, err)
	}

	return nil, fmt.Errorf("no proposal response found: %v", errors)
}

func parseResultFromTransactionAction(transactionAction *peer.TransactionAction) ([]byte, error) {
	actionPayload := &peer.ChaincodeActionPayload{}
	if err := proto.Unmarshal(transactionAction.GetPayload(), actionPayload); err != nil {
		return nil, fmt.Errorf("failed to deserialize chaincode action payload: %w", err)
	}

	responsePayload := &peer.ProposalResponsePayload{}
	if err := proto.Unmarshal(actionPayload.GetAction().GetProposalResponsePayload(), responsePayload); err != nil {
		return nil, fmt.Errorf("failed to deserialize proposal response payload: %w", err)
	}

	chaincodeAction := &peer.ChaincodeAction{}
	if err := proto.Unmarshal(responsePayload.GetExtension(), chaincodeAction); err != nil {
		return nil, fmt.Errorf("failed to deserialize chaincode action: %w", err)
	}

	return chaincodeAction.GetResponse().GetPayload(), nil
}
