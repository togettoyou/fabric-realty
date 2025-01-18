// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import "context"

type contextWithCancel func(parent context.Context) (context.Context, context.CancelFunc)

type contextFactory struct {
	ctx          context.Context
	evaluate     contextWithCancel
	endorse      contextWithCancel
	submit       contextWithCancel
	commitStatus contextWithCancel
}

func (factory *contextFactory) getOrDefault(supplier contextWithCancel) (context.Context, context.CancelFunc) {
	if supplier != nil {
		return supplier(factory.ctx)
	}
	return context.WithCancel(factory.ctx)
}

func (factory *contextFactory) Evaluate() (context.Context, context.CancelFunc) {
	return factory.getOrDefault(factory.evaluate)
}

func (factory *contextFactory) Endorse() (context.Context, context.CancelFunc) {
	return factory.getOrDefault(factory.endorse)
}

func (factory *contextFactory) Submit() (context.Context, context.CancelFunc) {
	return factory.getOrDefault(factory.submit)
}

func (factory *contextFactory) CommitStatus() (context.Context, context.CancelFunc) {
	return factory.getOrDefault(factory.commitStatus)
}
