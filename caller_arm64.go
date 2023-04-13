// Copyright 2023 The zapcl Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build go1.13 && arm64

package zapcl

import (
	"runtime"
)

func FuncForPC(pc uintptr) *runtime.Func {
	return runtime.FuncForPC(pc)
}

func CallersFrames(callers []uintptr) *runtime.Frames {
	return runtime.CallersFrames(callers)
}

func Caller(skip int) (pc uintptr, file string, line int, ok bool) {
	return runtime.Caller(skip)
}

// Callers is a drop-in replacement for runtime.Callers that uses frame
// pointers for fast and simple stack unwinding.
//
// Based by: https://github.com/golang/go/blob/go1.20/src/runtime/extern.go#L256-L264
//
//go:noinline
func Callers(skip int, pcs []uintptr) int {
	return runtime.Callers(skip, pcs)
}
