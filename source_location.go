// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	sourceLocationKey = "logging.googleapis.com/sourceLocation"
)

type sourceLocation struct {
	// File is the source file name. Depending on the runtime environment, this might be a simple name or a fully-qualified name.
	//
	// Optional.
	File string `json:"file"`

	// Line within the source file. 1-based; 0 indicates no line number available.
	//
	// Optional.
	Line string `json:"line"` // int64 format

	// Function is the Human-readable name of the function or method being invoked, with optional context such as the class or package name.
	//
	// This information may be used in contexts such as the logs viewer, where a file and line number are less meaningful.
	// The format can vary by language.
	// For example: dir/package.func.
	//
	// Optional.
	Function string `json:"function"`
}

// MarshalLogObject implements zapcore.ObjectMarshaller.MarshalLogObject.
func (l sourceLocation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("file", l.File)
	enc.AddString("line", l.Line)
	enc.AddString("function", l.Function)

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
		File:     file,
		Line:     strconv.Itoa(line),
		Function: function,
	}

	return loc
}

// SourceLocation adds the Cloud Logging "sourceLocation" field.
//
//	https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntrySourceLocation
func SourceLocation(pc uintptr, file string, line int, ok bool) zapcore.Field {
	return zap.Object(sourceLocationKey, newSource(pc, file, line, ok))
}
