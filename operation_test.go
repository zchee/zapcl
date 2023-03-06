// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestOperation(t *testing.T) {
	t.Parallel()

	op := &operation{
		LogEntryOperation: &logpb.LogEntryOperation{
			Id:       "id",
			Producer: "producer",
			First:    true,
			Last:     false,
		},
	}
	field := Operation("id", "producer", true, false)

	if diff := cmp.Diff(field, zap.Object(OperationKey, op),
		protocmp.Transform(),
	); diff != "" {
		t.Fatalf("(-want, +got)\n%s\n", diff)
	}
}

func TestOperationStart(t *testing.T) {
	t.Parallel()

	op := &operation{
		LogEntryOperation: &logpb.LogEntryOperation{
			Id:       "id",
			Producer: "producer",
			First:    true,
			Last:     false,
		},
	}
	field := OperationStart("id", "producer")

	if diff := cmp.Diff(field, zap.Object(OperationKey, op),
		protocmp.Transform(),
	); diff != "" {
		t.Fatalf("(-want, +got)\n%s\n", diff)
	}
}

func TestOperationCont(t *testing.T) {
	t.Parallel()

	op := &operation{
		LogEntryOperation: &logpb.LogEntryOperation{
			Id:       "id",
			Producer: "producer",
			First:    false,
			Last:     false,
		},
	}
	field := OperationCont("id", "producer")

	if diff := cmp.Diff(field, zap.Object(OperationKey, op),
		protocmp.Transform(),
	); diff != "" {
		t.Fatalf("(-want, +got)\n%s\n", diff)
	}
}

func TestOperationEnd(t *testing.T) {
	t.Parallel()

	op := &operation{
		LogEntryOperation: &logpb.LogEntryOperation{
			Id:       "id",
			Producer: "producer",
			First:    false,
			Last:     true,
		},
	}
	field := OperationEnd("id", "producer")

	if diff := cmp.Diff(field, zap.Object(OperationKey, op),
		protocmp.Transform(),
	); diff != "" {
		t.Fatalf("(-want, +got)\n%s\n", diff)
	}
}
