// Copyright 2023 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package json

import (
	"io"

	"github.com/bytedance/sonic/encoder"
)

// NewEncoder returns the new bytedance/sonic/encoder.StreamEncoder.
func NewEncoder(w io.Writer) *encoder.StreamEncoder {
	return encoder.NewStreamEncoder(w)
}
