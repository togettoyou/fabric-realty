// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"math"

	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-protos-go-apiv2/orderer"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func seekLargestBlockNumber() *orderer.SeekPosition {
	return &orderer.SeekPosition{
		Type: &orderer.SeekPosition_Specified{
			Specified: &orderer.SeekSpecified{
				Number: math.MaxUint64,
			},
		},
	}
}

type baseBlockEventsBuilder struct {
	eventsBuilder
	tlsCertificateHash []byte
}

func (builder *baseBlockEventsBuilder) payloadBytes() ([]byte, error) {
	channelHeader, err := builder.channelHeaderBytes()
	if err != nil {
		return nil, err
	}

	signatureHeader, err := builder.signatureHeaderBytes()
	if err != nil {
		return nil, err
	}

	data, err := builder.dataBytes()
	if err != nil {
		return nil, err
	}

	payload := &common.Payload{
		Header: &common.Header{
			ChannelHeader:   channelHeader,
			SignatureHeader: signatureHeader,
		},
		Data: data,
	}

	return proto.Marshal(payload)
}

func (builder *baseBlockEventsBuilder) channelHeaderBytes() ([]byte, error) {
	channelHeader := &common.ChannelHeader{
		Type:        int32(common.HeaderType_DELIVER_SEEK_INFO),
		Timestamp:   timestamppb.Now(),
		ChannelId:   builder.eventsBuilder.channelName,
		Epoch:       0,
		TlsCertHash: builder.tlsCertificateHash,
	}

	return proto.Marshal(channelHeader)
}

func (builder *baseBlockEventsBuilder) signatureHeaderBytes() ([]byte, error) {
	creator, err := builder.signingID.Creator()
	if err != nil {
		return nil, err
	}

	signatureHeader := &common.SignatureHeader{
		Creator: creator,
	}

	return proto.Marshal(signatureHeader)
}

func (builder *baseBlockEventsBuilder) dataBytes() ([]byte, error) {
	data := &orderer.SeekInfo{
		Start: builder.getStartPosition(),
		Stop:  seekLargestBlockNumber(),
	}

	return proto.Marshal(data)
}

type filteredBlockEventsBuilder struct {
	baseBlockEventsBuilder
}

func (builder *filteredBlockEventsBuilder) build() (*FilteredBlockEventsRequest, error) {
	payload, err := builder.payloadBytes()
	if err != nil {
		return nil, err
	}

	result := &FilteredBlockEventsRequest{
		baseBlockEventsRequest{
			client:    builder.client,
			signingID: builder.signingID,
			request: &common.Envelope{
				Payload: payload,
			},
		},
	}
	return result, nil
}

type blockEventsBuilder struct {
	baseBlockEventsBuilder
}

func (builder *blockEventsBuilder) build() (*BlockEventsRequest, error) {
	payload, err := builder.payloadBytes()
	if err != nil {
		return nil, err
	}

	result := &BlockEventsRequest{
		baseBlockEventsRequest{
			client:    builder.client,
			signingID: builder.signingID,
			request: &common.Envelope{
				Payload: payload,
			},
		},
	}
	return result, nil
}

type blockAndPrivateDataEventsBuilder struct {
	baseBlockEventsBuilder
}

func (builder *blockAndPrivateDataEventsBuilder) build() (*BlockAndPrivateDataEventsRequest, error) {
	payload, err := builder.payloadBytes()
	if err != nil {
		return nil, err
	}

	result := &BlockAndPrivateDataEventsRequest{
		baseBlockEventsRequest{
			client:    builder.client,
			signingID: builder.signingID,
			request: &common.Envelope{
				Payload: payload,
			},
		},
	}
	return result, nil
}
