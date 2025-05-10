// Package stacktrace provides utilities for capturing, representing, and
// formatting stack traces.
//
// The package allows you to capture the current call stack, inspect stack
// frames, and render stack traces as human-readable strings.
//
// Example:
//
//	stack := stacktrace.Capture(0, 10)
//	fmt.Println(stack)
//
// This captures up to 10 frames of the current stack and prints them in a
// readable format.
package stacktrace

import (
	"math"
	"runtime"
)

// DefaultDepth is the default maximum number of stack frames to capture.
const DefaultDepth = 64

// Capture captures the current stack trace, skipping the specified number of
// frames and limiting the depth of the trace. If depth is math.MaxInt, it
// captures the full trace.
//
// skip controls how many stack frames to skip (0 means start from the caller of Capture).
// depth limits the number of frames captured; use math.MaxInt for no limit.
//
// Returns a *Stack containing the captured frames.
func Capture(skip, depth int) *Stack {
	skip++ // we don't want current function ot get to trace

	var pcs []uintptr
	if depth == math.MaxInt {
		pcs = callersFull(skip)
	} else {
		pcs = callersFinite(skip, depth)
	}

	stack := NewStack(len(pcs))
	frames := runtime.CallersFrames(pcs)

	for {
		frame, next := frames.Next()

		stack.AddFrame(frame)

		if !next {
			break
		}
	}

	return stack
}

func callersFinite(skip, depth int) []uintptr {
	skip += 2 // we don't want current function and runtime.Callers ot get to trace

	pcs := make([]uintptr, depth)
	pcsLen := runtime.Callers(skip, pcs)

	return pcs[:pcsLen]
}

func callersFull(skip int) []uintptr {
	skip += 2 // we don't want current function and runtime.Callers ot get to trace

	pcs := make([]uintptr, DefaultDepth)
	pcsLen := runtime.Callers(skip, pcs)

	for pcsLen == len(pcs) {
		pcs = make([]uintptr, len(pcs)*2) //nolint:mnd
		pcsLen = runtime.Callers(skip, pcs)
	}

	return pcs[:pcsLen]
}
