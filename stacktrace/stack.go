package stacktrace

import (
	"iter"
	"runtime"
	"strconv"
	"strings"
)

type Stack struct {
	frames []runtime.Frame `exhaustruct:"optional"`
}

// AddCaller adds the caller's frame to the stack.
func (s *Stack) AddCaller() {
	pcs := make([]uintptr, 1)
	runtime.Callers(3, pcs) //nolint:mnd

	frame, _ := runtime.CallersFrames(pcs).Next()

	s.frames = append(s.frames, frame)
}

func (s *Stack) AddFrame(f runtime.Frame) {
	s.frames = append(s.frames, f)
}

// FramesIter returns an iterator over the frames in the stack.
//
// The iterator yields the frames starting from uppermost (the one added last).
func (s *Stack) FramesIter() iter.Seq2[int, runtime.Frame] {
	return func(yield func(int, runtime.Frame) bool) {
		start := len(s.frames) - 1
		for i := start; i >= 0; i-- {
			yield(start-i, s.frames[i])
		}
	}
}

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

// WriteFrameToBuffer writes a string representation of the given [runtime.Frame]
// to the provided [strings.Builder].
//
// The format is:
//
//		FunctionName
//	     FilePath:LineNumber
func WriteFrameToBuffer(f runtime.Frame, b *strings.Builder) {
	b.WriteString(f.Function)
	b.WriteRune('\n')
	b.WriteRune('\t')
	b.WriteString(f.File)
	b.WriteRune(':')
	b.WriteString(strconv.Itoa(f.Line))
}
