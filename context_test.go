// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

package zapcloudlogging

import (
	"strings"
	"testing"
)

func TestErrorReport(t *testing.T) {
	t.Parallel()

	got := ErrorReport(Caller(0)).Interface.(*reportContext)

	if gotFile, wantFile := got.ReportLocation.File, "zap-cloudlogging/context_test.go"; !strings.Contains(gotFile, wantFile) {
		t.Fatalf("except got.File contains %s but got %s", wantFile, gotFile)
	}
	if gotLine, wantLine := got.ReportLocation.Line, "14"; !strings.EqualFold(gotLine, wantLine) {
		t.Fatalf("except got.Line equal %s but got %s", wantLine, gotLine)
	}
	if gotFunc, wantFunc := got.ReportLocation.Function, "zap-cloudlogging.TestErrorReport"; !strings.Contains(gotFunc, wantFunc) {
		t.Fatalf("except got.Function contains %s but got %s", wantFunc, gotFunc)
	}
}

func TestNewReportContext(t *testing.T) {
	t.Parallel()

	got := newReportContext(Caller(0))

	if gotFile, wantFile := got.ReportLocation.File, "zap-cloudlogging/context_test.go"; !strings.Contains(gotFile, wantFile) {
		t.Fatalf("except got.File contains %s but got %s", wantFile, gotFile)
	}
	if gotLine, wantLine := got.ReportLocation.Line, "30"; !strings.EqualFold(gotLine, wantLine) {
		t.Fatalf("except got.Line equal %s but got %s", wantLine, gotLine)
	}
	if gotFunc, wantFunc := got.ReportLocation.Function, "zap-cloudlogging.TestNewReportContext"; !strings.Contains(gotFunc, wantFunc) {
		t.Fatalf("except got.Function contains %s but got %s", wantFunc, gotFunc)
	}
}
