package e_test

import (
	"errors"
	"os"
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
				in:       e.NewFrom("error", errors.New("wrapped")), //nolint:err113
				expected: "error: wrapped",
			},
			{
				name:     "wrapped error with fields",
				in:       e.NewFrom("error", e.New("wrapped", fields.F("f1", "v1")), fields.F("f2", "v2")),
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
		e2 := e.NewFrom("e2", e1)
		e3 := e.NewFrom("e3", e2)

		assert.Equal(t, "e3: e2: e1", e3.Error())
		assert.Same(t, e2, e3.Unwrap())
		assert.Same(t, e1, e2.Unwrap())
	})

	t.Run(".Wrap()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1", fields.F("f1", "v1"))
		e2 := e.New("e2", fields.F("f2", "v2"))

		assert.NotSame(t, e1, e1.Wrap(e2))
		assert.NotSame(t, e2, e1.Wrap(e2))
		assert.Equal(t, "e1 (f1=v1): e2 (f2=v2)", e1.Wrap(e2).Error())
		assert.Equal(t, "e1 (f1=v1) (f3=v3): e2 (f2=v2)", e1.Wrap(e2, fields.F("f3", "v3")).Error())
	})

	t.Run(".Unwrap()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1")
		e2 := e.NewFrom("e2", e1)
		e3 := e.NewFrom("e3", e2)

		assert.Same(t, e2, e3.Unwrap())
		assert.Same(t, e1, e2.Unwrap())
		assert.NoError(t, e1.Unwrap())
	})

	t.Run(".Is()", func(t *testing.T) {
		t.Parallel()

		var (
			e0       = errors.New("e0") //nolint:err113
			e1 error = e.NewFrom("e1", os.ErrNotExist)
			e2 error = e.From(e0)
			e3 error = e.NewFrom("e3", e1)
			e4       = e.From(e0)
		)

		assert.ErrorIs(t, e1, e1) //nolint:testifylint
		assert.NotErrorIs(t, e1, e0)

		assert.ErrorIs(t, e2, e0)

		assert.ErrorIs(t, e3, e1)
		assert.ErrorIs(t, e3, os.ErrNotExist)

		assert.NotErrorIs(t, e4, e2)
	})

	t.Run(".As()", func(t *testing.T) {
		t.Parallel()

		var (
			e0 error = &myErr{"e0"}
			e1       = e.New("e1")
			e2 error = e.From(e0)
			e3 error = e.NewFrom("e3", e2)
			e4       = e1.Wrap(e3)
		)

		var target *myErr

		assert.False(t, errors.As(e1, &target))

		assert.ErrorAs(t, e2, &target)
		assert.ErrorAs(t, e3, &target)
		assert.ErrorAs(t, e4, &target)
	})

	t.Run(".Reason()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1", fields.F("f1", "v1"))
		e2 := e.NewFrom("e2", e1, fields.F("f2", "v2"))

		assert.Equal(t, "e1", e1.Reason())
		assert.Equal(t, "e2", e2.Reason())
	})

	t.Run(".Fields()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1", fields.F("f1", "v1"))
		e2 := e.NewFrom("e2", e1, fields.F("f2", "v2"))

		assert.Equal(t, fields.List{fields.F("f1", "v1")}, e1.Fields())
		assert.Equal(t, fields.List{fields.F("f2", "v2")}, e2.Fields())
	})

	t.Run(".Clone()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1", fields.F("f1", "v1"))
		e2 := e1.Clone()

		assert.NotSame(t, e1, e2)
		assert.Equal(t, e1.Error(), e2.Error())
		assert.NotSame(t, e1.Fields(), e2.Fields())
		assert.ElementsMatch(t, e1.Fields(), e2.Fields())
	})

	t.Run(".WithFields()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1", fields.F("f1", "v1"))
		e2 := e1.WithFields(fields.F("f2", "v2"))

		assert.NotSame(t, e1, e2)
		assert.NotSame(t, e1.Fields(), e2.Fields())
		assert.Len(t, e1.Fields(), 1)
		assert.Len(t, e2.Fields(), 2)
	})

	t.Run(".WithField()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1", fields.F("f1", "v1"))
		e2 := e1.WithField("f2", "v2")

		assert.NotSame(t, e1, e2)
		assert.NotSame(t, e1.Fields(), e2.Fields())
		assert.Len(t, e1.Fields(), 1)
		assert.Len(t, e2.Fields(), 2)
	})
}

type myErr struct {
	err string
}

func (err *myErr) Error() string {
	return err.err
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
			in:       []any{e.NewFrom("e2", e.New("e1")), e.NewFrom("e4", e.New("e3"))},
			expected: "e2: e1: e4: e3",
		},
	}

	for _, tc := range tt {
		assert.Equal(t, tc.expected, e.Wrap(tc.in...).Error(), tc.name)
	}
}
