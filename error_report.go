// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcl

import (
	"go/build"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
)

const (
	contextKey = "context"
)

// reportLocation is the source code location information associated with the log entry
// for the purpose of reporting an error, if any.
type reportLocation struct {
	*logpb.LogEntrySourceLocation
}

// MarshalLogObject implements zapcore.ObjectMarshaller.MarshalLogObject.
func (l reportLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("filePath", l.GetFile())
	enc.AddInt64("lineNumber", l.GetLine())
	enc.AddString("functionName", l.GetFunction())

	return nil
}

// reportContext is the context information attached to a log for reporting errors.
type reportContext struct {
	ReportLocation *reportLocation `json:"reportLocation"`
}

// MarshalLogObject implements zapcore.ObjectMarshaller.MarshalLogObject.
func (c reportContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	return enc.AddObject("reportLocation", c.ReportLocation)
}

func newReportContext(pc uintptr, file string, line int, ok bool) *reportContext {
	if !ok {
		return nil
	}

	var function string
	if fn := FuncForPC(pc); fn != nil {
		function = strings.TrimPrefix(fn.Name(), filepath.Join(build.Default.GOPATH, "src")+"/")
	}
	ctx := &reportContext{
		ReportLocation: &reportLocation{
			LogEntrySourceLocation: &logpb.LogEntrySourceLocation{
				File:     file,
				Line:     int64(line),
				Function: function,
			},
		},
	}

	return ctx
}

// ErrorReport adds the Cloud Logging "context" field for getting the log line reported as error.
//
// https://cloud.google.com/error-reporting/docs/formatting-error-messages
func ErrorReport(pc uintptr, file string, line int, ok bool) zap.Field {
	return zap.Object(contextKey, newReportContext(pc, file, line, ok))
}
