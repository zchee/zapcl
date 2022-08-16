// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"strings"
	"testing"
)

func TestSourceLocation(t *testing.T) {
	t.Parallel()

	got := SourceLocation(Caller(0)).Interface.(*sourceLocation)

	if gotFile, wantFile := got.File, "zap-cloudlogging/source_location_test.go"; !strings.Contains(gotFile, wantFile) {
		t.Fatalf("except got.File contains %s but got %s", wantFile, gotFile)
	}
	if gotLine, wantLine := got.Line, int64(14); gotLine != wantLine {
		t.Fatalf("except got.Line equal %d but got %d", wantLine, gotLine)
	}
	if gotFunc, wantFunc := got.Function, "zap-cloudlogging.TestSourceLocation"; !strings.Contains(gotFunc, wantFunc) {
		t.Fatalf("except got.Function contains %s but got %s", wantFunc, gotFunc)
	}
}

func TestNewSource(t *testing.T) {
	t.Parallel()

	got := newSource(Caller(0))

	if gotFile, wantFile := got.File, "zap-cloudlogging/source_location_test.go"; !strings.Contains(gotFile, wantFile) {
		t.Fatalf("except got.File contains %s but got %s", wantFile, gotFile)
	}
	if gotLine, wantLine := got.Line, int64(30); gotLine != wantLine {
		t.Fatalf("except got.Line equal %d but got %d", wantLine, gotLine)
	}
	if gotFunc, wantFunc := got.Function, "zap-cloudlogging.TestNewSource"; !strings.Contains(gotFunc, wantFunc) {
		t.Fatalf("except got.Function contains %s but got %s", wantFunc, gotFunc)
	}
}
