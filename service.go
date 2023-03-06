// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcl

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	serviceContextKey = "serviceContext"
)

// ServiceContext adds the service information adding the log line.
// It is a required field if an error needs to be reported.
//
// ServiceContext: https://cloud.google.com/error-reporting/reference/rest/v1beta1/ServiceContext
func ServiceContext(name string) zap.Field {
	return zap.Object(serviceContextKey, zapcore.ObjectMarshalerFunc(func(enc zapcore.ObjectEncoder) error {
		enc.AddString("service", name)

		return nil
	}))
}
