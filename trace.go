// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"context"
	"fmt"

	"cloud.google.com/go/compute/metadata"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	traceKey        = "logging.googleapis.com/trace"
	spanKey         = "logging.googleapis.com/spanId"
	traceSampledKey = "logging.googleapis.com/trace_sampled"
)

// TraceField adds the correct Cloud Logging "trace", "span", "trace_sampled" fields from ctx.
//
//	https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
func TraceField(ctx context.Context) []zapcore.Field {
	projectID, err := metadata.ProjectID()
	if err != nil {
		panic(err)
	}

	spanCtx := trace.SpanContextFromContext(ctx)

	return []zapcore.Field{
		zap.String(traceKey, fmt.Sprintf("projects/%s/traces/%s", projectID, spanCtx.TraceID().String())),
		zap.String(spanKey, spanCtx.SpanID().String()),
		zap.Bool(traceSampledKey, spanCtx.IsSampled()),
	}
}
