// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcl

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/zchee/zapcl/pkg/monitoredresource"
)

const (
	// TraceKey is the resource name of the trace associated with the log entry if any. For more information.
	//
	// trace field:
	// - https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
	// - https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#FIELDS.trace
	TraceKey = "logging.googleapis.com/trace"

	// SpanKey is the span ID within the trace associated with the log entry. For more information.
	//
	// spanId field:
	// - https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
	// - https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#FIELDS.span_id
	SpanKey = "logging.googleapis.com/spanId"

	// TraceSampledKey is the value of this field must be either true or false. For more information.
	//
	// trace_sampled field:
	// - https://cloud.google.com/logging/docs/structured-logging#special-payload-fields
	// - https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#FIELDS.trace_sampled
	TraceSampledKey = "logging.googleapis.com/trace_sampled"
)

// TraceField adds the correct Cloud Logging "trace", "span", "trace_sampled" fields from ctx.
//
// https://cloud.google.com/logging/docs/agent/logging/configuration#special-fields
func TraceField(ctx context.Context) (fields []zapcore.Field, record bool) {
	span := trace.SpanFromContext(ctx)

	if record = span.IsRecording(); !record {
		return nil, record
	}

	return []zapcore.Field{
		zap.String(TraceKey, fmt.Sprintf("projects/%s/traces/%s", monitoredresource.ResourceDetector.ProjectID(), span.SpanContext().TraceID().String())),
		zap.String(SpanKey, span.SpanContext().SpanID().String()),
		zap.Bool(TraceSampledKey, span.SpanContext().IsSampled()),
	}, true
}
