// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/zap"
)

func TestLabel(t *testing.T) {
	t.Parallel()

	field := Label("key", "value")

	if diff := cmp.Diff(field, zap.String("labels.key", "value")); diff != "" {
		t.Fatalf("(-want, +got)\n%s\n", diff)
	}
}

func TestLabels(t *testing.T) {
	t.Parallel()

	field := Labels(
		"hello", "world",
		"hi", "universe",
	)

	labels := new(labelMap)
	for key, val := range map[string]string{"hello": "world", "hi": "universe"} {
		labels.Add(key, val)
	}

	if diff := cmp.Diff(field, zap.Object(LabelsKey, labels),
		cmpopts.IgnoreUnexported(labelMap{}),
	); diff != "" {
		t.Fatalf("(-want, +got)\n%s\n", diff)
	}
}
