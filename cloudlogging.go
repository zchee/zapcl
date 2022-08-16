// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	cloudlogging "cloud.google.com/go/logging"
	"go.uber.org/zap/zapcore"
)

var levelToSeverity = map[zapcore.Level]cloudlogging.Severity{
	zapcore.DebugLevel:  cloudlogging.Debug,
	zapcore.InfoLevel:   cloudlogging.Info,
	zapcore.WarnLevel:   cloudlogging.Warning,
	zapcore.ErrorLevel:  cloudlogging.Error,
	zapcore.DPanicLevel: cloudlogging.Critical,
	zapcore.PanicLevel:  cloudlogging.Alert,
	zapcore.FatalLevel:  cloudlogging.Emergency,
}

type logger struct {
	zapcore.LevelEnabler
}

var _ zapcore.Core = (*logger)(nil)

// With adds structured context to the Core.
func (l *logger) With(fields []zapcore.Field) zapcore.Core {
	return nil
}

// Check determines whether the supplied Entry should be logged (using the
// embedded LevelEnabler and possibly some extra logic). If the entry
// should be logged, the Core adds itself to the CheckedEntry and returns
// the result.
func (l *logger) Check(entry zapcore.Entry, centry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return nil
}

// Write serializes the Entry and any Fields supplied at the log site and
// writes them to their destination.
func (l *logger) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	return nil
}

// Sync flushes buffered logs (if any).
func (l *logger) Sync() error {
	return nil
}
