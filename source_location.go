// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
)

const (
	sourceLocationKey = "logging.googleapis.com/sourceLocation"
)

type sourceLocation struct {
	*logpb.LogEntrySourceLocation
}

// MarshalLogObject implements zapcore.ObjectMarshaller.MarshalLogObject.
func (l sourceLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("file", l.GetFile())
	enc.AddInt64("line", l.GetLine())
	enc.AddString("function", l.GetFunction())

	return nil
}

func newSource(pc uintptr, file string, line int, ok bool) *sourceLocation {
	if !ok {
		return nil
	}

	var function string
	if fn := FuncForPC(pc); fn != nil {
		function = fn.Name()
	}

	loc := &sourceLocation{
		LogEntrySourceLocation: &logpb.LogEntrySourceLocation{
			File:     file,
			Line:     int64(line),
			Function: function,
		},
	}

	return loc
}

// SourceLocation adds the Cloud Logging "sourceLocation" field.
//
//	https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntrySourceLocation
func SourceLocation(pc uintptr, file string, line int, ok bool) zapcore.Field {
	return zap.Object(sourceLocationKey, newSource(pc, file, line, ok))
}
