package stacktrace

import (
	"iter"
	"math"
	"runtime"
	"strconv"
	"strings"
)

// Stack represents a captured call stack consisting of program counter frames.
type Stack struct {
	frames []runtime.Frame
}

// Frames returns an iterator over the stack frames.
// The iterator yields (index, frame) pairs for each frame in the stack.
func (s *Stack) Frames() iter.Seq2[int, runtime.Frame] {
	return func(yield func(int, runtime.Frame) bool) {
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
	b := strings.Builder{}
	b.Grow(len(s.frames) * 100)

	for i, f := range s.frames {
		if i > 0 {
			b.WriteRune('\n')
		}

		WriteFrameToBuffer(f, &b)
	}

	return b.String()
}

// WriteFrameToBuffer formats a single runtime.Frame and writes it to the
// provided buffer. The frame is formatted as "function\n\tfile:line".
func WriteFrameToBuffer(f runtime.Frame, b *strings.Builder) {
	b.WriteString(f.Function)
	b.WriteString("\n\t")
	b.WriteString(f.File)
	b.WriteRune(':')
	b.WriteString(strconv.Itoa(f.Line))
}

// DefaultDepth is the default maximum number of stack frames to capture.
const DefaultDepth = 64

// Capture captures the current goroutine's call stack. The skip argument
// specifies the number of frames to skip before recording (0 means start at
// caller). The depth argument specifies the maximum number of frames to capture
// (use math.MaxInt for unlimited).
func Capture(skip, depth int) *Stack {
	skip++ // we don't want current function to get to trace

	var pcs []uintptr
	if depth == math.MaxInt {
		pcs = callersFull(skip)
	} else {
		pcs = callersFinite(skip, depth)
	}

	stack := &Stack{frames: make([]runtime.Frame, 0, len(pcs))}
	frames := runtime.CallersFrames(pcs)

	for {
		frame, next := frames.Next()

		stack.frames = append(stack.frames, frame)

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
