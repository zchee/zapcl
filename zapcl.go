// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcl

import (
	"errors"
	"fmt"
	"sort"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sys/unix"
	logtypepb "google.golang.org/genproto/googleapis/logging/type"

	"github.com/zchee/zapcl/internal/json"
	"github.com/zchee/zapcl/pkg/monitoredresource"
)

const (
	timeKey       = "time" // https://cloud.google.com/logging/docs/agent/logging/configuration#timestamp-processing
	levelKey      = "severity"
	nameKey       = "logger"
	callerKey     = "caller"
	messageKey    = "message"
	stacktraceKey = "stacktrace"
)

var levelToSeverity = map[zapcore.Level]logtypepb.LogSeverity{
	zapcore.DebugLevel:  logtypepb.LogSeverity_DEBUG,
	zapcore.InfoLevel:   logtypepb.LogSeverity_INFO,
	zapcore.WarnLevel:   logtypepb.LogSeverity_WARNING,
	zapcore.ErrorLevel:  logtypepb.LogSeverity_ERROR,
	zapcore.DPanicLevel: logtypepb.LogSeverity_CRITICAL,
	zapcore.PanicLevel:  logtypepb.LogSeverity_ALERT,
	zapcore.FatalLevel:  logtypepb.LogSeverity_EMERGENCY,
}

func levelEncoder(lvl zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(levelToSeverity[lvl].Enum().String())
}

// NewEncoderConfig returns the logging configuration.
func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:             timeKey,
		LevelKey:            levelKey,
		NameKey:             nameKey,
		CallerKey:           callerKey,
		MessageKey:          messageKey,
		StacktraceKey:       stacktraceKey,
		LineEnding:          zapcore.DefaultLineEnding,
		EncodeLevel:         levelEncoder,
		EncodeTime:          zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration:      zapcore.SecondsDurationEncoder,
		EncodeCaller:        zapcore.ShortCallerEncoder,
		NewReflectedEncoder: json.NewEncoder,
	}
}

type nopWriteSyncer struct{}

// Write implements zapcore.WriteSyncer.
func (nopWriteSyncer) Write([]byte) (int, error) { return 0, nil }

// Sync implements zapcore.WriteSyncer.
func (nopWriteSyncer) Sync() error { return nil }

var _ zapcore.WriteSyncer = nopWriteSyncer{}

// core represents a zapcor.core that is Cloud Logging integration for Zap logger.
type core struct {
	zapcore.LevelEnabler

	enc        zapcore.Encoder
	ws         zapcore.WriteSyncer
	initFields map[string]any
	fields     []zapcore.Field
}

var _ zapcore.Core = (*core)(nil)

func newCore(ws zapcore.WriteSyncer, enab zapcore.LevelEnabler, opts ...Option) *core {
	encoderConfig := NewEncoderConfig()

	core := &core{
		LevelEnabler: enab,
		enc:          zapcore.NewJSONEncoder(encoderConfig),
		ws:           ws,
	}
	for _, opt := range opts {
		opt.apply(core)
	}

	res := monitoredresource.Detect()
	core.fields = []zapcore.Field{
		zap.String(res.Type, res.LogID),
		zap.Inline(res),
	}

	// handling initFields option
	if len(core.initFields) > 0 {
		fs := make([]zapcore.Field, 0, len(core.initFields))
		keys := make([]string, 0, len(core.initFields))
		for k := range core.initFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fs = append(fs, zap.Any(k, core.initFields[k]))
		}
		core.fields = append(core.fields, fs...)
	}

	return core
}

// NewCore creates a Core that writes logs to a WriteSyncer.
func NewCore(ws zapcore.WriteSyncer, enab zapcore.LevelEnabler, opts ...Option) zapcore.Core {
	core := newCore(ws, enab, opts...)

	return zapcore.NewCore(core.enc, core.ws, core.LevelEnabler)
}

// WrapCore wraps or replaces the Logger's underlying zapcore.Core.
func WrapCore(opts ...Option) zap.Option {
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		core := newCore(nopWriteSyncer{}, c, opts...)

		return zapcore.NewCore(core.enc, core.ws, core.LevelEnabler)
	})
}

// Write serializes the Entry and any Fields supplied at the log site and
// writes them to their destination.
//
// Write implements zapcore.Core.
func (c *core) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	for _, field := range c.fields {
		field.AddTo(c.enc)
	}

	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return fmt.Errorf("could not encode entry: %w", err)
	}

	_, err = c.ws.Write(buf.Bytes())
	buf.Free()
	if err != nil {
		return fmt.Errorf("could not write buf: %w", err)
	}

	if ent.Level > zapcore.ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		c.Sync() //nolint:errcheck
	}

	return nil
}

// Check determines whether the supplied Entry should be logged (using the
// embedded LevelEnabler and possibly some extra logic). If the entry
// should be logged, the Core adds itself to the CheckedEntry and returns
// the result.
//
// Check implements zapcore.Core.
func (c *core) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}

	return ce
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

func (c *core) clone() *core {
	newCore := &core{
		LevelEnabler: c.LevelEnabler,
		fields:       make([]zapcore.Field, len(c.fields)),
		enc:          c.enc.Clone(),
		ws:           c.ws,
	}
	copy(newCore.fields, c.fields)

	return newCore
}

// With adds structured context to the Core.
//
// With implements zapcore.Core.
func (c *core) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	addFields(clone.enc, fields)

	return clone
}

// Sync flushes buffered logs if any.
//
// Sync implements zapcore.Core.
func (c *core) Sync() error {
	if err := c.ws.Sync(); err != nil {
		if !knownSyncError(err) {
			return fmt.Errorf("sync logger: %w", err)
		}
	}

	return nil
}

var knownSyncErrors = []error{
	// sync /dev/stdout: invalid argument
	unix.EINVAL,
	// sync /dev/stdout: not supported
	unix.ENOTSUP,
	// sync /dev/stdout: inappropriate ioctl for device
	unix.ENOTTY,
	// sync /dev/stdout: bad file descriptor
	unix.EBADF,
}

// knownSyncError returns true if the given error is one of the known
// non-actionable errors returned by Sync on Linux and macOS.
//
// This code was borrowed from https://github.com/open-telemetry/opentelemetry-collector/blob/v0.74.0/exporter/loggingexporter/known_sync_error.go#L25-L46.
func knownSyncError(err error) bool {
	for _, syncError := range knownSyncErrors {
		if errors.Is(err, syncError) {
			return true
		}
	}

	return false
}
