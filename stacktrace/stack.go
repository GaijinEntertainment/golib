package stacktrace

import (
	"iter"
	"runtime"
	"strconv"
	"strings"
)

// Stack represents a captured stack trace as a sequence of frames.
//
// Use the Capture function to obtain a Stack, and the String method to render it.
//
// Example:
//
//	stack := stacktrace.Capture(0, 10)
//	fmt.Println(stack)
type Stack struct {
	frames []runtime.Frame
}

// NewStack returns a new Stack with preallocated space for the given initial size.
func NewStack(initialSize int) *Stack {
	return &Stack{
		frames: make([]runtime.Frame, 0, initialSize),
	}
}

// AddCaller adds the current caller's frame to the stack.
func (s *Stack) AddCaller() {
	pcs := make([]uintptr, 1)
	runtime.Callers(3, pcs) //nolint:mnd

	frame, _ := runtime.CallersFrames(pcs).Next()

	s.AddFrame(frame)
}

// AddFrame appends the given runtime.Frame to the stack.
func (s *Stack) AddFrame(f runtime.Frame) {
	s.frames = append(s.frames, f)
}

// Len returns the number of frames contained in the stack.
func (s *Stack) Len() int {
	return len(s.frames)
}

// FramesIter returns an iterator over the frames in the stack.
//
// The iterator yields the frames starting from the most recent (the one added last).
func (s *Stack) FramesIter() iter.Seq2[int, runtime.Frame] {
	return func(yield func(int, runtime.Frame) bool) {
		start := len(s.frames) - 1
		for i := start; i >= 0; i-- {
			if !yield(start-i, s.frames[i]) {
				break
			}
		}
	}
}

// String returns a multi-line string representation of the stack trace.
// Each frame is shown with its function name and file:line.
func (s *Stack) String() string {
	b := &strings.Builder{}

	for _, f := range s.FramesIter() {
		if b.Len() > 1 {
			b.WriteRune('\n')
		}

		WriteFrameToBuffer(f, b)
	}

	return b.String()
}

// WriteFrameToBuffer writes a string representation of the given runtime.Frame
// to the provided strings.Builder.
//
// The format is:
//
//	FunctionName
//	\tFilePath:LineNumber
func WriteFrameToBuffer(f runtime.Frame, b *strings.Builder) {
	b.WriteString(f.Function)
	b.WriteRune('\n')
	b.WriteRune('\t')
	b.WriteString(f.File)
	b.WriteRune(':')
	b.WriteString(strconv.Itoa(f.Line))
}
