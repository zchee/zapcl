// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"context"

	"go.uber.org/zap"
)

// ctxKey is how we find Loggers in a context.Context.
type ctxKey struct{}

// FromContext returns a *zap.Logger from ctx.
func FromContext(ctx context.Context) *zap.Logger {
	return ctx.Value(ctxKey{}).(*zap.Logger)
}

// NewContext returns a new Context, derived from ctx, which carries the provided *zap.Logger.
func NewContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}
