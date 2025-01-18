// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"

	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"google.golang.org/protobuf/proto"
)

type chaincodeEventsBuilder struct {
	eventsBuilder
	chaincodeName string
}

func (builder *chaincodeEventsBuilder) build() (*ChaincodeEventsRequest, error) {
	signedRequest, err := builder.newSignedChaincodeEventsRequestProto()
	if err != nil {
		return nil, err
	}

	result := &ChaincodeEventsRequest{
		client:        builder.client,
		signingID:     builder.signingID,
		signedRequest: signedRequest,
	}
	return result, nil
}

func (builder *chaincodeEventsBuilder) newSignedChaincodeEventsRequestProto() (*gateway.SignedChaincodeEventsRequest, error) {
	request, err := builder.newChaincodeEventsRequestProto()
	if err != nil {
		return nil, err
	}

	requestBytes, err := proto.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize chaincode events request: %w", err)
	}

	signedRequest := &gateway.SignedChaincodeEventsRequest{
		Request: requestBytes,
	}
	return signedRequest, nil
}

func (builder *chaincodeEventsBuilder) newChaincodeEventsRequestProto() (*gateway.ChaincodeEventsRequest, error) {
	creator, err := builder.signingID.Creator()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize identity: %w", err)
	}

	request := &gateway.ChaincodeEventsRequest{
		ChannelId:          builder.channelName,
		Identity:           creator,
		ChaincodeId:        builder.chaincodeName,
		StartPosition:      builder.getStartPosition(),
		AfterTransactionId: builder.afterTransactionID,
	}
	return request, nil
}
