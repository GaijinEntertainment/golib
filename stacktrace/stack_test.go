package stacktrace_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/stacktrace"
)

func takeOddCallers(s *stacktrace.Stack, n int) {
	if n == 1 {
		return
	}

	if n%2 == 1 {
		s.AddCaller()
	}

	takeOddCallers(s, n-1)
}

func TestStack(t *testing.T) {
	t.Parallel()

	t.Run("AddCaller", func(t *testing.T) {
		t.Parallel()

		s := stacktrace.NewStack(0)
		s.AddCaller()

		takeOddCallers(s, 5)

		require.Equal(t, 3, s.Len())

		frames := make([]runtime.Frame, 3)
		for i, f := range s.FramesIter() {
			frames[i] = f
		}

		assert.Equal(t, "dev.gaijin.team/go/golib/stacktrace_test.takeOddCallers", frames[0].Function)
		assert.True(t, strings.HasSuffix(frames[0].File, "/stacktrace/stack_test.go"))
		assert.Equal(t, 23, frames[0].Line)

		assert.Equal(t, "dev.gaijin.team/go/golib/stacktrace_test.TestStack.func1", frames[1].Function)
		assert.True(t, strings.HasSuffix(frames[1].File, "/stacktrace/stack_test.go"))
		assert.Equal(t, 35, frames[1].Line)

		assert.Equal(t, "testing.tRunner", frames[2].Function)
	})

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		s := stacktrace.NewStack(0)
		s.AddCaller()

		takeOddCallers(s, 5)

		lines := strings.Split(s.String(), "\n")

		require.Len(t, lines, 6)

		assert.Equal(t, "dev.gaijin.team/go/golib/stacktrace_test.takeOddCallers", lines[0])
		assert.Equal(t, "dev.gaijin.team/go/golib/stacktrace_test.TestStack.func2", lines[2])

		assert.True(t, strings.HasPrefix(lines[1], "\t"))
		assert.True(t, strings.HasSuffix(lines[1], "stacktrace/stack_test.go:23"))
		assert.True(t, strings.HasPrefix(lines[3], "\t"))
		assert.True(t, strings.HasSuffix(lines[3], "stacktrace/stack_test.go:61"))
		assert.True(t, strings.HasPrefix(lines[5], "\t"))
	})
}
