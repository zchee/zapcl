// Copyright 2023 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build !go1.21 && amd64
// +build !go1.21,amd64

package json

import (
	"io"

	"github.com/bytedance/sonic"
	"go.uber.org/zap/zapcore"
)

// NewEncoder returns the new bytedance/sonic/encoder.StreamEncoder.
func NewEncoder(w io.Writer) zapcore.ReflectedEncoder {
	enc := sonic.ConfigFastest.NewEncoder(w)
	enc.SetEscapeHTML(false)

	return enc
}
