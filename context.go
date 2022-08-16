// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	contextKey = "context"
)

// reportLocation is the source code location information associated with the log entry
// for the purpose of reporting an error, if any.
type reportLocation struct {
	File     string `json:"filePath"`
	Line     string `json:"lineNumber"`
	Function string `json:"functionName"`
}

// MarshalLogObject implements zapcore.ObjectMarshaller.MarshalLogObject.
func (l reportLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("filePath", l.File)
	enc.AddString("lineNumber", l.Line)
	enc.AddString("functionName", l.Function)

	return nil
}

// reportContext is the context information attached to a log for reporting errors.
type reportContext struct {
	ReportLocation reportLocation `json:"reportLocation"`
}

// MarshalLogObject implements zapcore.ObjectMarshaller.MarshalLogObject.
func (c reportContext) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddObject("reportLocation", c.ReportLocation)

	return nil
}

func newReportContext(pc uintptr, file string, line int, ok bool) *reportContext {
	if !ok {
		return nil
	}

	var function string
	if fn := FuncForPC(pc); fn != nil {
		function = fn.Name()
	}
	ctx := &reportContext{
		ReportLocation: reportLocation{
			File:     file,
			Line:     strconv.Itoa(line),
			Function: function,
		},
	}

	return ctx
}

// ErrorReport adds the Cloud Logging "context" field for getting the log line reported as error.
//
//	https://cloud.google.com/error-reporting/docs/formatting-error-messages
func ErrorReport(pc uintptr, file string, line int, ok bool) zap.Field {
	return zap.Object(contextKey, newReportContext(pc, file, line, ok))
}
