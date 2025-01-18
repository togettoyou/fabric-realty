// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"

	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/gateway"
	"github.com/hyperledger/fabric-protos-go-apiv2/peer"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type gatewayClient struct {
	grpcGatewayClient gateway.GatewayClient
	grpcDeliverClient peer.DeliverClient
	contexts          *contextFactory
}

func (client *gatewayClient) Endorse(in *gateway.EndorseRequest, opts ...grpc.CallOption) (*gateway.EndorseResponse, error) {
	ctx, cancel := client.contexts.Endorse()
	defer cancel()
	return client.EndorseWithContext(ctx, in, opts...)
}

func (client *gatewayClient) EndorseWithContext(ctx context.Context, in *gateway.EndorseRequest, opts ...grpc.CallOption) (*gateway.EndorseResponse, error) {
	response, err := client.grpcGatewayClient.Endorse(ctx, in, opts...)
	if err != nil {
		txErr := newTransactionError(err, in.GetTransactionId())
		return nil, &EndorseError{txErr}
	}

	return response, nil
}

func (client *gatewayClient) Submit(in *gateway.SubmitRequest, opts ...grpc.CallOption) (*gateway.SubmitResponse, error) {
	ctx, cancel := client.contexts.Submit()
	defer cancel()
	return client.SubmitWithContext(ctx, in, opts...)
}

func (client *gatewayClient) SubmitWithContext(ctx context.Context, in *gateway.SubmitRequest, opts ...grpc.CallOption) (*gateway.SubmitResponse, error) {
	response, err := client.grpcGatewayClient.Submit(ctx, in, opts...)
	if err != nil {
		txErr := newTransactionError(err, in.GetTransactionId())
		return nil, &SubmitError{txErr}
	}

	return response, nil
}

func (client *gatewayClient) CommitStatus(in *gateway.SignedCommitStatusRequest, opts ...grpc.CallOption) (*gateway.CommitStatusResponse, error) {
	ctx, cancel := client.contexts.CommitStatus()
	defer cancel()
	return client.CommitStatusWithContext(ctx, in, opts...)
}

func (client *gatewayClient) CommitStatusWithContext(ctx context.Context, in *gateway.SignedCommitStatusRequest, opts ...grpc.CallOption) (*gateway.CommitStatusResponse, error) {
	response, err := client.grpcGatewayClient.CommitStatus(ctx, in, opts...)
	if err != nil {
		transactionID := getTransactionIDFromSignedCommitStatusRequest(in)
		txErr := newTransactionError(err, transactionID)
		return nil, &CommitStatusError{txErr}
	}

	return response, nil
}

func getTransactionIDFromSignedCommitStatusRequest(in *gateway.SignedCommitStatusRequest) string {
	request := &gateway.CommitStatusRequest{}
	err := proto.Unmarshal(in.GetRequest(), request)
	if err != nil {
		return "?"
	}
	return request.GetTransactionId()
}

func (client *gatewayClient) Evaluate(in *gateway.EvaluateRequest, opts ...grpc.CallOption) (*gateway.EvaluateResponse, error) {
	ctx, cancel := client.contexts.Evaluate()
	defer cancel()
	return client.EvaluateWithContext(ctx, in, opts...)
}

func (client *gatewayClient) EvaluateWithContext(ctx context.Context, in *gateway.EvaluateRequest, opts ...grpc.CallOption) (*gateway.EvaluateResponse, error) {
	return client.grpcGatewayClient.Evaluate(ctx, in, opts...)
}

func (client *gatewayClient) ChaincodeEvents(ctx context.Context, in *gateway.SignedChaincodeEventsRequest, opts ...grpc.CallOption) (gateway.Gateway_ChaincodeEventsClient, error) {
	return client.grpcGatewayClient.ChaincodeEvents(ctx, in, opts...)
}

func (client *gatewayClient) BlockEvents(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (peer.Deliver_DeliverClient, error) {
	deliverClient, err := client.grpcDeliverClient.Deliver(ctx, opts...)
	if err != nil {
		return nil, err
	}

	if err := deliverClient.Send(in); err != nil {
		return nil, err
	}

	return deliverClient, nil
}

func (client *gatewayClient) FilteredBlockEvents(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (peer.Deliver_DeliverFilteredClient, error) {
	deliverClient, err := client.grpcDeliverClient.DeliverFiltered(ctx, opts...)
	if err != nil {
		return nil, err
	}

	if err := deliverClient.Send(in); err != nil {
		return nil, err
	}

	return deliverClient, nil
}

func (client *gatewayClient) BlockAndPrivateDataEvents(ctx context.Context, in *common.Envelope, opts ...grpc.CallOption) (peer.Deliver_DeliverWithPrivateDataClient, error) {
	deliverClient, err := client.grpcDeliverClient.DeliverWithPrivateData(ctx, opts...)
	if err != nil {
		return nil, err
	}

	if err := deliverClient.Send(in); err != nil {
		return nil, err
	}

	return deliverClient, nil
}
