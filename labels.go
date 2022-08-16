// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	labelsKey = "logging.googleapis.com/labels"
)

// Label adds the Cloud Logging "labels" field from key and val.
//
// Cloud Logging truncates label keys that exceed 512 B and label values that exceed 64 KB upon their associated log entry being written.
// The truncation is indicated by an ellipsis at the end of the character string.
//
//	https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#FIELDS.labels
func Label(key, val string) zapcore.Field {
	return zap.String("labels."+key, val)
}

// labelMap map of labels.
type labelMap struct {
	m sync.Map
}

// Add adds the val for a key.
func (l *labelMap) Add(key, val string) {
	l.m.Store(key, val)
}

// Delete deletes the value for a key.
func (l *labelMap) Delete(key string) {
	l.m.Delete(key)
}

// MarshalLogObject implements zapcore.ObjectMarshaler.
func (l *labelMap) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	l.m.Range(func(key, val any) bool {
		enc.AddString(key.(string), val.(string))

		return true
	})

	return nil
}

// Label adds the Cloud Logging "labels" field from keyvals.
func Labels(keyvals ...string) zapcore.Field {
	if len(keyvals)%2 != 0 {
		panic("keyval length should be powers of two")
	}

	fields := new(labelMap)
	for i := 0; i < len(keyvals)/2; i++ {
		key := keyvals[i]
		val := keyvals[i+1]
		fields.Add(key, val)
	}

	return zap.Object(labelsKey, fields)
}
