package stacktrace_test

import (
	"math"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/stacktrace"
)

func TestFrame(t *testing.T) {
	t.Parallel()

	t.Run("NewFrame copies necessary fields", func(t *testing.T) {
		t.Parallel()

		// Capture a real frame
		var pcs [1]uintptr

		n := runtime.Callers(1, pcs[:])
		require.Equal(t, 1, n, "should capture one frame")

		runtimeFrame, _ := runtime.CallersFrames(pcs[:]).Next()

		frame := stacktrace.NewFrame(runtimeFrame)

		assert.Equal(t, runtimeFrame.PC, frame.PC, "PC should match")
		assert.Equal(t, runtimeFrame.File, frame.File, "File should match")
		assert.Equal(t, runtimeFrame.Line, frame.Line, "Line should match")
		assert.Equal(t, runtimeFrame.Function, frame.Function, "Function should match")
	})

	t.Run("FullPath formats correctly", func(t *testing.T) {
		t.Parallel()

		frame := stacktrace.Frame{
			PC:       1,
			File:     "/path/to/file.go",
			Line:     42,
			Function: "package.Function",
		}

		fullPath := frame.FullPath()
		assert.Equal(t, "/path/to/file.go:42", fullPath)
	})

	t.Run("FullPath handles zero PC", func(t *testing.T) {
		t.Parallel()

		frame := stacktrace.Frame{
			PC:       0,
			File:     "/path/to/file.go",
			Line:     42,
			Function: "package.Function",
		}

		fullPath := frame.FullPath()
		assert.Equal(t, "undefined", fullPath)
	})

	t.Run("String formats correctly", func(t *testing.T) {
		t.Parallel()

		frame := stacktrace.Frame{
			PC:       1,
			File:     "/Users/test/main.go",
			Line:     100,
			Function: "main.run",
		}

		assert.Equal(t, "main.run\n\t/Users/test/main.go:100", frame.String())
	})

	t.Run("String handles zero PC", func(t *testing.T) {
		t.Parallel()

		frame := stacktrace.Frame{
			PC:       0,
			File:     "/path/to/file.go",
			Line:     42,
			Function: "package.Function",
		}

		assert.Equal(t, "undefined", frame.String())
	})
}

func recursionA(i, depth int) *stacktrace.Stack {
	if i == 0 {
		return stacktrace.CaptureStack(0, depth)
	}

	fn := recursionB
	if i%2 == 0 {
		fn = recursionA
	}

	return fn(i-1, depth)
}

func recursionB(i, depth int) *stacktrace.Stack {
	if i == 0 {
		return stacktrace.CaptureStack(0, depth)
	}

	fn := recursionB
	if i%2 == 0 {
		fn = recursionA
	}

	return fn(i-1, depth)
}

func TestCaptureStack(t *testing.T) {
	t.Parallel()

	assert.Empty(t, new(stacktrace.Stack).String(), "empty stack should produce empty string")

	t.Run("shallow stack", func(t *testing.T) {
		t.Parallel()

		// Capture a very shallow stack with only 2 frames
		s := stacktrace.CaptureStack(0, 2)
		require.Equal(t, 2, s.Len())

		// Test through the .String function to also check stringification validity.
		str := s.String()
		require.NotEmpty(t, str)

		// Verify exact structure: should have exactly one newline separator between frames
		lines := strings.Split(str, "\n")
		require.Len(t, lines, 4, "should have 4 lines: func1, location1, func2, location2")

		// First frame: function name (the current test function)
		require.Contains(t, lines[0], "TestCaptureStack", "first frame should be test function")
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

func TestCaptureCaller(t *testing.T) {
	t.Parallel()

	t.Run("captures direct caller", func(t *testing.T) {
		t.Parallel()

		frame := stacktrace.CaptureCaller(0)
		require.NotEmpty(t, frame.Function, "should capture function name")
		require.Contains(t, frame.Function, "TestCaptureCaller", "should contain test function name")
		require.NotEmpty(t, frame.File, "should capture file path")
		require.Contains(t, frame.File, "stacktrace_test.go", "should be from test file")
		require.Positive(t, frame.Line, "should have positive line number")
	})

	t.Run("skip parameter works", func(t *testing.T) {
		t.Parallel()

		// CaptureStack from this test directly with different skip values
		frame0 := stacktrace.CaptureCaller(0)
		require.Contains(t, frame0.Function, "TestCaptureCaller", "skip=0 should capture test function")

		// Use helper to test skip=1
		helperFunc := func() stacktrace.Frame {
			return stacktrace.CaptureCaller(1) // skip=1 skips helper, captures test
		}

		frame1 := helperFunc()
		require.Contains(t, frame1.Function, "TestCaptureCaller", "skip=1 should skip helper and capture test")
	})

	t.Run("returns empty frame when no caller", func(t *testing.T) {
		t.Parallel()

		// Very large skip should result in empty frame
		frame := stacktrace.CaptureCaller(1000)
		require.Empty(t, frame.Function, "should return empty frame for too-large skip")
	})
}
