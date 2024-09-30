package e_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/GaijinEntertainment/golib/e"
	"github.com/GaijinEntertainment/golib/fields"
)

func TestErr(t *testing.T) {
	t.Parallel()

	t.Run(".Error()", func(t *testing.T) {
		t.Parallel()

		tt := []struct {
			name     string
			in       *e.Err
			expected string
		}{
			{
				name:     "nil",
				in:       nil,
				expected: "nil",
			},
			{
				name:     "empty",
				in:       &e.Err{},
				expected: "nil",
			},
			{
				name:     "simple",
				in:       e.New("reason"),
				expected: "reason",
			},
			{
				name:     "with fields",
				in:       e.New("reason", fields.F("key", "value")),
				expected: "reason (key=value)",
			},
			{
				name:     "wrapped error",
				in:       e.NewFrom(e.New("wrapped"), "error"),
				expected: "error: wrapped",
			},
			{
				name:     "wrapped error with fields",
				in:       e.NewFrom(e.New("wrapped", fields.F("f1", "v1")), "error", fields.F("f2", "v2")),
				expected: "error (f2=v2): wrapped (f1=v1)",
			},
			{
				name:     "nil error in the middle",
				in:       e.Wrap("e1", nil, "e2"),
				expected: "e1: <nil>(<nil>): e2",
			},
		}

		for _, tc := range tt {
			assert.Equal(t, tc.expected, tc.in.Error(), tc.name)
		}
	})

	t.Run("Unwrap()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1")
		e2 := e.NewFrom(e1, "e2")
		e3 := e.NewFrom(e2, "e3")

		assert.Equal(t, "e3: e2: e1", e3.Error())
		assert.Same(t, e2, e3.Unwrap())
		assert.Same(t, e1, e2.Unwrap())
	})

	t.Run(".Wrap()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1", fields.F("f1", "v1"))
		e2 := e.New("e2")

		assert.NotSame(t, e1, e1.Wrap(e2))
		assert.NotSame(t, e2, e1.Wrap(e2))
		assert.Equal(t, "e1 (f1=v1): e2", e1.Wrap(e2).Error())
		assert.Equal(t, "e1 (f1=v1) (f3=v3): e2", e1.Wrap(e2, fields.F("f3", "v3")).Error())
	})
}

func TestWrap(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		in       []any
		expected string
	}{
		{
			name:     "empty",
			in:       nil,
			expected: "nil",
		},
		{
			name:     "string",
			in:       []any{"error"},
			expected: "error",
		},
		{
			name:     "error",
			in:       []any{e.New("error")},
			expected: "error",
		},
		{
			name:     "e error",
			in:       []any{e.New("e1")},
			expected: "e1",
		},
		{
			name:     "int32",
			in:       []any{int32(42)},
			expected: "int32(42)",
		},
		{
			name:     "multiple error",
			in:       []any{"e1", "e2", "e3"},
			expected: "e1: e2: e3",
		},
		{
			name:     "multiple errors with fields",
			in:       []any{e.New("e1", fields.F("f1", "v1")), "e2", "e3"},
			expected: "e1 (f1=v1): e2: e3",
		},
		{
			name:     "multiple wrapped errors",
			in:       []any{e.NewFrom(e.New("e1"), "e2"), e.NewFrom(e.New("e3"), "e4")},
			expected: "e2: e1: e4: e3",
		},
	}

	for _, tc := range tt {
		assert.Equal(t, tc.expected, e.Wrap(tc.in...).Error(), tc.name)
	}
}
