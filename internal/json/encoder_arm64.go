// Copyright 2023 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package json

import (
	"io"

	gojson "github.com/goccy/go-json"
)

// NewEncoder returns the github.com/goccy/go-json NewEncoder.
func NewEncoder(w io.Writer) *gojson.Encoder {
	return gojson.NewEncoder(w)
}
