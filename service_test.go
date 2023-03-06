// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestServiceContext(t *testing.T) {
	t.Parallel()

	field := zapcore.Field{
		Key:  "serviceContext",
		Type: zapcore.ObjectMarshalerType,
	}

	if diff := cmp.Diff(field, ServiceContext("test service name"),
		cmpopts.IgnoreFields(zap.Field{}, "Interface"),
	); diff != "" {
		t.Fatalf("(-want, +got)\n%s\n", diff)
	}
}
