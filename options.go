// Copyright 2023 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcl

import (
	"go.uber.org/zap/zapcore"
)

// Option configures a core.
type Option interface {
	apply(*core)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*core)

func (f optionFunc) apply(c *core) {
	f(c)
}

// WithInitialFields configures the zap InitialFields.
func WithInitialFields(fields map[string]any) Option {
	return optionFunc(func(c *core) {
		c.initFields = fields
	})
}

// WithWriteSyncer configures the zapcore.WriteSyncer.
func WithWriteSyncer(ws zapcore.WriteSyncer) Option {
	return optionFunc(func(c *core) {
		c.ws = ws
	})
}
