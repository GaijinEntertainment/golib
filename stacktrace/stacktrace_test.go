package stacktrace_test

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/stacktrace"
)

func recursionA(i, depth int) *stacktrace.Stack {
	if i == 0 {
		return stacktrace.Capture(0, depth)
	}

	fn := recursionB
	if i%2 == 0 {
		fn = recursionA
	}

	return fn(i-1, depth)
}

func recursionB(i, depth int) *stacktrace.Stack {
	if i == 0 {
		return stacktrace.Capture(0, depth)
	}

	fn := recursionB
	if i%2 == 0 {
		fn = recursionA
	}

	return fn(i-1, depth)
}

func TestCapture(t *testing.T) {
	t.Parallel()

	t.Run("shallow stack", func(t *testing.T) {
		t.Parallel()

		// Capture a very shallow stack with only 2 frames
		s := stacktrace.Capture(0, 2)
		require.Equal(t, 2, s.Len())

		str := s.String()
		require.NotEmpty(t, str)

		// Verify exact structure: should have exactly one newline separator between frames
		lines := strings.Split(str, "\n")
		require.Len(t, lines, 4, "should have 4 lines: func1, location1, func2, location2")

		// First frame: function name (the current test function)
		require.Contains(t, lines[0], "TestCapture", "first frame should be test function")
		require.Contains(t, lines[0], "stacktrace_test", "first frame should be in stacktrace_test package")

		// First frame: location (starts with tab)
		require.True(t, strings.HasPrefix(lines[1], "\t"), "location should start with tab")
		require.Contains(t, lines[1], "stacktrace_test.go:", "should contain file:line")

		// Second frame: function name
		require.NotEmpty(t, lines[2], "second frame function should not be empty")

		// Second frame: location (starts with tab)
		require.True(t, strings.HasPrefix(lines[3], "\t"), "location should start with tab")
		require.Contains(t, lines[3], ".go:", "should contain file:line")
	})

	t.Run("finite depth", func(t *testing.T) {
		t.Parallel()

		s := recursionA(33, 42)

		// it is expected stack to be 4 elements bigger since tests
		// itself are also saturating callstack
		require.Equal(t, 37, s.Len())

		s = recursionA(42, 32)

		// in this case it is exact match since callstack is bigger that required depth
		require.Equal(t, 32, s.Len())
	})

	t.Run("full depth", func(t *testing.T) {
		t.Parallel()

		s := recursionA(33, math.MaxInt)
		require.Equal(t, 37, s.Len())

		s = recursionA(146, math.MaxInt)
		require.Equal(t, 150, s.Len())

		// Test Frames() iterator
		count := 0
		for idx, frame := range s.Frames() {
			require.Equal(t, count, idx, "frame index should match iteration count")
			require.NotEmpty(t, frame.Function, "frame should have a function name")
			require.NotEmpty(t, frame.File, "frame should have a file path")
			require.Positive(t, frame.Line, "frame should have a positive line number")

			count = idx + 1
		}

		require.Equal(t, s.Len(), count, "iterator should yield same number of frames as Len()")
	})
}
