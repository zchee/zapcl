// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"cloud.google.com/go/logging"
	"go.uber.org/zap/zapcore"
)

var levelToSeverity = map[zapcore.Level]logging.Severity{
	zapcore.DebugLevel:  logging.Debug,
	zapcore.InfoLevel:   logging.Info,
	zapcore.WarnLevel:   logging.Warning,
	zapcore.ErrorLevel:  logging.Error,
	zapcore.DPanicLevel: logging.Critical,
	zapcore.PanicLevel:  logging.Alert,
	zapcore.FatalLevel:  logging.Emergency,
}
