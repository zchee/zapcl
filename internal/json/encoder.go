// Copyright 2023 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build arm64 || (go1.21 && amd64)
// +build arm64 go1.21,amd64

package json

import (
	"io"

	gojson "github.com/goccy/go-json"
	"go.uber.org/zap/zapcore"
)

// NewEncoder returns the github.com/goccy/go-json NewEncoder.
func NewEncoder(w io.Writer) zapcore.ReflectedEncoder {
	enc := gojson.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc
}
