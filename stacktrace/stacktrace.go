package stacktrace

import (
	"bytes"
	"iter"
	"math"
	"runtime"
	"strconv"
)

// Frame represents a single stack frame containing program counter, file path,
// line number, and function name.
//
// This is a simplified representation of runtime.Frame that is safe to store and
// use after the stack trace has been captured.
type Frame struct {
	// PC is the program counter for this frame.
	PC uintptr
	// File is the full path to the source file.
	File string
	// Line is the line number in the source file.
	Line int
	// Function is the name of the function.
	Function string
}

// NewFrame creates a Frame from a runtime.Frame, copying the essential fields.
func NewFrame(f runtime.Frame) Frame {
	return Frame{
		PC:       f.PC,
		File:     f.File,
		Line:     f.Line,
		Function: f.Function,
	}
}

// FullPath returns the full file path and line number formatted as
// "/path/to/file:line". Returns "undefined" if the frame has a zero PC.
func (f Frame) FullPath() string {
	b := &bytes.Buffer{}
	f.WriteFullPath(b)
	return b.String()
}

// WriteFullPath writes the full file path and line number of the frame to the
// provided buffer in the format "/path/to/file:line". If the PC is zero, it writes
// "undefined".
func (f Frame) WriteFullPath(b *bytes.Buffer) {
	if f.PC == 0 {
		b.WriteString("undefined")
		return
	}

	b.WriteString(f.File)
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(f.Line))
}

// String returns a formatted representation of the frame including function
// name and source location. Returns "undefined" if the frame has a zero PC.
func (f Frame) String() string {
	b := &bytes.Buffer{}
	f.Write(b)
	return b.String()
}

// Write writes the function name and source location of the frame to the
// provided buffer in the format "function\n\t/path/to/file:line". If the PC is
// zero, it writes "undefined".
func (f Frame) Write(b *bytes.Buffer) {
	if f.PC == 0 {
		b.WriteString("undefined")
		return
	}

	b.WriteString(f.Function)
	b.WriteString("\n\t")
	b.WriteString(f.File)
	b.WriteByte(':')
	b.WriteString(strconv.Itoa(f.Line))
}

// Stack represents a captured call stack consisting of program counter frames.
type Stack struct {
	frames []Frame
}

// Frames returns an iterator over the stack frames.
// The iterator yields (index, frame) pairs for each frame in the stack.
func (s *Stack) Frames() iter.Seq2[int, Frame] {
	return func(yield func(int, Frame) bool) {
		for i, frame := range s.frames {
			if !yield(i, frame) {
				break
			}
		}
	}
}

// Len returns the number of frames in the stack trace.
func (s *Stack) Len() int {
	return len(s.frames)
}

// String formats the stack trace as a multi-line string with function names and
// source locations. Each frame is formatted as "function\n\tfile:line".
func (s *Stack) String() string {
	if len(s.frames) == 0 {
		return ""
	}

	// Estimate capacity: ~100 bytes per frame on average
	b := bytes.NewBuffer(make([]byte, 0, len(s.frames)*100)) //nolint:mnd

	for i, f := range s.frames {
		if i > 0 {
			b.WriteRune('\n')
		}

		f.Write(b)
	}

	return b.String()
}

// DefaultDepth is the default maximum number of stack frames to capture.
const DefaultDepth = 64

// CaptureCaller captures a single caller frame from the current goroutine's call
// stack. The skip argument specifies the number of frames to skip before
// recording (0 means the direct caller of CaptureCaller).
//
// This is more efficient than CaptureStack when only the immediate caller information
// is needed, as it only captures and processes a single frame.
func CaptureCaller(skip int) Frame {
	skip += 2 // we don't want current function and runtime.Callers to get to trace

	var pcs [1]uintptr

	n := runtime.Callers(skip, pcs[:])

	if n == 0 {
		return Frame{} //nolint:exhaustruct
	}

	runtimeFrame, _ := runtime.CallersFrames(pcs[:n]).Next()

	return NewFrame(runtimeFrame)
}

// CaptureStack captures the current goroutine's call stack. The skip argument
// specifies the number of frames to skip before recording (0 means start at
// caller). The depth argument specifies the maximum number of frames to capture
// (use math.MaxInt for unlimited).
func CaptureStack(skip, depth int) *Stack {
	skip++ // we don't want current function to get to trace

	var pcs []uintptr
	if depth == math.MaxInt {
		pcs = callersFull(skip)
	} else {
		pcs = callersFinite(skip, depth)
	}

	stack := &Stack{frames: make([]Frame, 0, len(pcs))}
	frames := runtime.CallersFrames(pcs)

	for {
		frame, next := frames.Next()

		stack.frames = append(stack.frames, NewFrame(frame))

		if !next {
			break
		}
	}

	return stack
}

func callersFinite(skip, depth int) []uintptr {
	skip += 2 // we don't want current function and runtime.Callers to get to trace

	pcs := make([]uintptr, depth)
	pcsLen := runtime.Callers(skip, pcs)

	return pcs[:pcsLen]
}

func callersFull(skip int) []uintptr {
	skip += 2 // we don't want current function and runtime.Callers to get to trace

	// Start with a reasonable default that handles most stacks efficiently.
	pcs := make([]uintptr, DefaultDepth)
	pcsLen := runtime.Callers(skip, pcs)

	// If the stack is deeper than our initial buffer, grow and recapture.
	// Note:
	// runtime.Callers always captures from the beginning, so we must recapture the
	// entire stack on each iteration. This is inherent to the API.
	for pcsLen == len(pcs) {
		pcs = make([]uintptr, len(pcs)*2) //nolint:mnd
		pcsLen = runtime.Callers(skip, pcs)
	}

	return pcs[:pcsLen]
}
