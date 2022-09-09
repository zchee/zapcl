// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcl

import (
	"strings"
	"testing"
)

func TestSourceLocation(t *testing.T) {
	t.Parallel()

	got := SourceLocation(Caller(0)).Interface.(*sourceLocation)

	if gotFile, wantFile := got.File, "zap-cloudlogging/source_location_test.go"; !strings.Contains(gotFile, wantFile) {
		t.Errorf("except contains got %s in %s", gotFile, wantFile)
	}
	if gotLine, wantLine := got.Line, int64(14); gotLine != wantLine {
		t.Errorf("except equal got %d in %d", gotLine, wantLine)
	}
	if gotFunc, wantFunc := got.Function, "zap-cloudlogging.TestSourceLocation"; !strings.Contains(gotFunc, wantFunc) {
		t.Errorf("except contains got %s in %s", gotFunc, wantFunc)
	}
}

func TestNewSource(t *testing.T) {
	t.Parallel()

	got := newSource(Caller(0))

	if gotFile, wantFile := got.File, "zap-cloudlogging/source_location_test.go"; !strings.Contains(gotFile, wantFile) {
		t.Fatalf("except contains got %s in %s", gotFile, wantFile)
	}
	if gotLine, wantLine := got.Line, int64(30); gotLine != wantLine {
		t.Errorf("except equal got %d in %d", gotLine, wantLine)
	}
	if gotFunc, wantFunc := got.Function, "zap-cloudlogging.TestNewSource"; !strings.Contains(gotFunc, wantFunc) {
		t.Errorf("except contains got %s in %s", gotFunc, wantFunc)
	}
}
