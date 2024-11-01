package stacktrace_test

import (
	"math"
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
	})
}
