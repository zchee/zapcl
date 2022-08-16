// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	operationKey = "logging.googleapis.com/operation"
)

// operation is the payload of Cloud Logging operation field.
type operation struct {
	// ID is an arbitrary operation identifier. Log entries with the same identifier are assumed to be part of the same operation.
	//
	// Optional.
	ID string `json:"id"`

	// Producer is an arbitrary producer identifier. The combination of id and producer must be globally unique.
	//
	// Examples for producer: "MyDivision.MyBigCompany.com", "github.com/MyProject/MyApplication".
	//
	// Optional.
	Producer string `json:"producer"`

	// First set this to True if this is the first log entry in the operation.
	//
	// Optional.
	First bool `json:"first"`

	// Last set this to True if this is the last log entry in the operation.
	//
	// Optional.
	Last bool `json:"last"`
}

var _ zapcore.ObjectMarshaler = (*operation)(nil)

// MarshalLogObject implements zapcore.ObjectMarshaller.MarshalLogObject.
func (op operation) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", op.ID)
	enc.AddString("producer", op.Producer)
	enc.AddBool("first", op.First)
	enc.AddBool("last", op.Last)

	return nil
}

// Operation adds the Cloud Logging "operation" fields from args.
//
//	https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#LogEntryOperation
func Operation(id, producer string, first, last bool) zapcore.Field {
	op := &operation{
		ID:       id,
		Producer: producer,
		First:    first,
		Last:     last,
	}

	return zap.Object(operationKey, op)
}

// OperationStart is a convenience function for `Operation`.
//
// It should be called for the first operation log.
func OperationStart(id, producer string) zapcore.Field {
	return Operation(id, producer, true, false)
}

// OperationCont is a convenience function for `Operation`.
//
// It should be called for any non-start/end operation log.
func OperationCont(id, producer string) zapcore.Field {
	return Operation(id, producer, false, false)
}

// OperationEnd is a convenience function for `Operation`.
//
// It should be called for the last operation lo
func OperationEnd(id, producer string) zapcore.Field {
	return Operation(id, producer, false, true)
}
