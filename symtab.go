// Copyright 2022 The zap-cloudlogging Authors
// SPDX-License-Identifier: BSD-3-Clause

//go:build go1.13

package zapcloudlogging

import (
	"runtime"
	"unsafe"
)

// Token from:
//
//	https://github.com/golang/go/blob/go1.19/src/runtime/symtab.go#L799-L802
type funcInfo struct {
	Func *runtime.Func  // *_func
	_    unsafe.Pointer // datap *moduledata
}

// https://github.com/golang/go/blob/go1.19/src/runtime/symtab.go#L828
//
//go:linkname findfunc runtime.findfunc
//go:noescape
func findfunc(pc uintptr) funcInfo

// FuncForPC is a drop-in replacement for runtime.FuncForPC.
//
//go:nosplit
func FuncForPC(pc uintptr) *runtime.Func {
	return findfunc(pc).Func
}

// frame is the information returned by Frames for each call frame.
//
// Token from:
//
//	https://github.com/golang/go/blob/go1.19/src/runtime/symtab.go#L26-L61
type frame struct {
	// PC is the program counter for the location in this frame.
	// For a frame that calls another frame, this will be the
	// program counter of a call instruction. Because of inlining,
	// multiple frames may have the same PC value, but different
	// symbolic information.
	PC uintptr

	// Func is the Func value of this call frame. This may be nil
	// for non-Go code or fully inlined functions.
	Func *runtime.Func

	// Function is the package path-qualified function name of
	// this call frame. If non-empty, this string uniquely
	// identifies a single function in the program.
	// This may be the empty string if not known.
	// If Func is not nil then Function == Func.Name().
	Function string

	// File and Line are the file name and line number of the
	// location in this frame. For non-leaf frames, this will be
	// the location of a call. These may be the empty string and
	// zero, respectively, if not known.
	File string
	Line int

	// Entry point program counter for the function; may be zero
	// if not known. If Func is not nil then Entry ==
	// Func.Entry().
	Entry uintptr

	// The runtime's internal view of the function. This field
	// is set (funcInfo.valid() returns true) only for Go functions,
	// not for C functions.
	funcInfo funcInfo
}

// frames may be used to get function/file/line information for a
// slice of PC values returned by Callers.
//
// Token from:
//
//	https://github.com/golang/go/blob/go1.19/src/runtime/symtab.go#L16-L23
type frames struct {
	// callers is a slice of PCs that have not yet been expanded to frames.
	callers []uintptr

	// frames is a slice of Frames that have yet to be returned.
	frames     []frame
	frameStore [2]frame
}

// CallersFrames is a drop-in replacement for runtime.CallersFrames.
//
// CallersFrames takes a slice of PC values returned by Callers and
// prepares to return function/file/line information.
// Do not change the slice until you are done with the Frames.
//
// Token from:
//
//	https://github.com/golang/go/blob/go1.19/src/runtime/symtab.go#L66
func CallersFrames(callers []uintptr) *frames {
	f := &frames{callers: callers}
	f.frames = f.frameStore[:0]
	return f
}

// https://github.com/golang/go/blob/go1.19/src/runtime/symtab.go#L81
//
//go:linkname next runtime.(*Frames).Next
//go:noescape
func next(*frames) (frame frame, more bool)

// Next returns a Frame representing the next call frame in the slice
// of PC values. If it has already returned all call frames, Next
// returns a zero Frame.
//
// The more result indicates whether the next call to Next will return
// a valid Frame. It does not necessarily indicate whether this call
// returned one.
//
// See the Frames example for idiomatic usage.
//
//go:nosplit
func (ci *frames) Next() (frame frame, more bool) {
	return next(ci)
}

// Caller is a drop-in replacement for runtime.Caller.
//
// Token from:
//
//	https://github.com/golang/go/blob/go1.19/src/runtime/extern.go#L217-L225
func Caller(skip int) (pc uintptr, file string, line int, ok bool) {
	rpc := make([]uintptr, 1)
	n := callers(skip+1, rpc[:])
	if n < 1 {
		return
	}
	frame, _ := CallersFrames(rpc).Next()

	return frame.PC, frame.File, frame.Line, frame.PC != 0
}

// Callers is a drop-in replacement for runtime.Callers that uses frame
// pointers for fast and simple stack unwinding.
//
// Based by:
//
//	https://github.com/golang/go/blob/go1.19/src/runtime/extern.go#L240-L248
//
//go:noinline
func Callers(skip int, pcs []uintptr) int {
	return callers(skip+1, pcs)
}

//go:noinline
//go:nosplit
func callers(skip int, pcs []uintptr) int {
	fp := uintptr(unsafe.Pointer(&skip)) - 16

	i := 0
	for i < len(pcs) {
		pc := deref(fp + 8)
		if skip == 0 {
			pcs[i] = pc
			i++
		} else {
			skip--
		}
		fp = deref(fp)
		if fp == 0 {
			break
		}
	}

	return i
}

//go:nosplit
func deref(addr uintptr) uintptr {
	return uintptr(**(**unsafe.Pointer)(unsafe.Pointer(&addr)))
}
