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

// ChaincodeEventsRequest delivers events emitted by transaction functions in a specific chaincode.
type ChaincodeEventsRequest struct {
	client        *gatewayClient
	signingID     *signingIdentity
	signedRequest *gateway.SignedChaincodeEventsRequest
}

// Bytes of the serialized chaincode events request.
func (events *ChaincodeEventsRequest) Bytes() ([]byte, error) {
	requestBytes, err := proto.Marshal(events.signedRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall SignedChaincodeEventsRequest protobuf: %w", err)
	}

	return requestBytes, nil
}

// Digest of the chaincode events request. This is used to generate a digital signature.
func (events *ChaincodeEventsRequest) Digest() []byte {
	return events.signingID.Hash(events.signedRequest.GetRequest())
}

// Events returns a channel from which chaincode events can be read.
func (events *ChaincodeEventsRequest) Events(ctx context.Context, opts ...grpc.CallOption) (<-chan *ChaincodeEvent, error) {
	if err := events.sign(); err != nil {
		return nil, err
	}

	eventsClient, err := events.client.ChaincodeEvents(ctx, events.signedRequest, opts...)
	if err != nil {
		return nil, err
	}

	results := make(chan *ChaincodeEvent)
	go func() {
		defer close(results)

		for {
			response, err := eventsClient.Recv()
			if err != nil {
				return
			}

			deliverChaincodeEvents(response, results)
		}
	}()

	return results, nil
}

func (events *ChaincodeEventsRequest) sign() error {
	if events.isSigned() {
		return nil
	}

	digest := events.Digest()
	signature, err := events.signingID.Sign(digest)
	if err != nil {
		return err
	}

	events.setSignature(signature)

	return nil
}

func (events *ChaincodeEventsRequest) isSigned() bool {
	return len(events.signedRequest.Signature) > 0
}

func (events *ChaincodeEventsRequest) setSignature(signature []byte) {
	events.signedRequest.Signature = signature
}

// ChaincodeEvent emitted by a transaction function.
type ChaincodeEvent struct {
	BlockNumber   uint64
	TransactionID string
	ChaincodeName string
	EventName     string
	Payload       []byte
}

func deliverChaincodeEvents(response *gateway.ChaincodeEventsResponse, send chan<- *ChaincodeEvent) {
	for _, event := range response.GetEvents() {
		send <- &ChaincodeEvent{
			BlockNumber:   response.GetBlockNumber(),
			TransactionID: event.GetTxId(),
			ChaincodeName: event.GetChaincodeId(),
			EventName:     event.GetEventName(),
			Payload:       event.GetPayload(),
		}
	}
}
