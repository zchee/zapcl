// Copyright 2022 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcl

import (
	"strings"
	"testing"
)

func TestErrorReport(t *testing.T) {
	t.Parallel()

	got := ErrorReport(Caller(0)).Interface.(*reportContext)

	if gotFile, wantFile := got.ReportLocation.File, "zapcl/error_report_test.go"; !strings.Contains(gotFile, wantFile) {
		t.Errorf("except contains got %s in %s", gotFile, wantFile)
	}
	if gotLine, wantLine := got.ReportLocation.Line, int64(14); gotLine != wantLine {
		t.Errorf("except equal got %d in %d", gotLine, wantLine)
	}
	if gotFunc, wantFunc := got.ReportLocation.Function, "zapcl.TestErrorReport"; !strings.Contains(gotFunc, wantFunc) {
		t.Errorf("except contains got %s in %s", gotFunc, wantFunc)
	}
}

func TestNewReportContext(t *testing.T) {
	t.Parallel()

	got := newReportContext(Caller(0))

	if gotFile, wantFile := got.ReportLocation.File, "zapcl/error_report_test.go"; !strings.Contains(gotFile, wantFile) {
		t.Fatalf("except contains got %s in %s", gotFile, wantFile)
	}
	if gotLine, wantLine := got.ReportLocation.Line, int64(30); gotLine != wantLine {
		t.Fatalf("except equal got %d in %d", gotLine, wantLine)
	}
	if gotFunc, wantFunc := got.ReportLocation.Function, "zapcl.TestNewReportContext"; !strings.Contains(gotFunc, wantFunc) {
		t.Errorf("except contains got %s in %s", gotFunc, wantFunc)
	}
}
