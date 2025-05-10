package e_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/fields"
)

func TestErr(t *testing.T) {
	t.Parallel()

	t.Run("general logic", func(t *testing.T) {
		t.Parallel()

		// general usecases
		var (
			errNotFound       = errors.New("not found") //nolint:err113
			errInvalidRequest = e.New("invalid request")
			errInvalidID      = errInvalidRequest.Wrap(errors.New("invalid id"))     //nolint:err113
			errInvalidAction  = errInvalidRequest.Wrap(errors.New("invalid action")) //nolint:err113
		)

		e1 := errInvalidAction.WithField("foo", "bar")
		e2 := errInvalidAction.WithField("foo", "bar")

		assert.NotErrorIs(t, errInvalidRequest, errNotFound)

		// they're both children of errInvalidRequest, therefore should pass error.Is
		// check.
		assert.ErrorIs(t, errInvalidID, errInvalidRequest)
		assert.ErrorIs(t, errInvalidAction, errInvalidRequest)

		// but they should fail if checked against each other
		assert.NotErrorIs(t, errInvalidID, errInvalidAction)
		assert.NotErrorIs(t, errInvalidAction, errInvalidID)

		// and errInvalidRequest should vail if tested against its children
		assert.NotErrorIs(t, errInvalidRequest, errInvalidID)
		assert.NotErrorIs(t, errInvalidRequest, errInvalidAction)

		// two instances with fields are not the same errors
		assert.NotSame(t, e1, e2)

		// but they're both instances of same error, therefore pass checks against
		// parents
		assert.ErrorIs(t, e1, errInvalidRequest)
		assert.ErrorIs(t, e2, errInvalidRequest)
		assert.ErrorIs(t, e1, errInvalidAction)
		assert.ErrorIs(t, e2, errInvalidAction)

		// but not against sibling of parent
		assert.NotErrorIs(t, e1, errInvalidID)
		assert.NotErrorIs(t, e2, errInvalidID)
	})

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
				expected: "(*e.Err)(nil)",
			},
			{
				name:     "nil wrapping fields",
				in:       (*e.Err)(nil).WithField("foo", "bar"),
				expected: "(*e.Err)(nil) (foo=bar)",
			},
			{
				name:     "empty",
				in:       &e.Err{},
				expected: "(*e.Err)(empty)",
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
				name: "wrapped error",
				//nolint:err113
				in:       e.NewFrom("error", errors.New("wrapped")),
				expected: "error: wrapped",
			},
			{
				name:     "wrapped error with fields",
				in:       e.NewFrom("error", e.New("wrapped", fields.F("f1", "v1")), fields.F("f2", "v2")),
				expected: "error (f2=v2): wrapped (f1=v1)",
			},
			{
				name:     "nil error in the middle",
				in:       e.New("e1").Wrap((*e.Err)(nil)).Wrap(e.New("e2")),
				expected: "e1: (*e.Err)(nil): e2",
			},
			{
				name:     "from external error with fields",
				in:       e.From(errors.New("error"), fields.F("key", "value")), //nolint:err113
				expected: "error (key=value)",
			},
			{
				name:     "from nil error",
				in:       e.From(nil),
				expected: "error(nil)",
			},
			{
				name:     "from nil with fields",
				in:       e.From(nil, fields.F("key", "value")),
				expected: "error(nil) (key=value)",
			},
			{
				name:     "empty with fields",
				in:       (&e.Err{}).WithField("foo", "bar"),
				expected: "(*e.Err)(empty) (foo=bar)",
			},
		}

		for _, tc := range tt {
			assert.Equal(t, tc.expected, tc.in.Error(), tc.name)
		}
	})

	t.Run(".Wrap()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1", fields.F("f1", "v1"))
		e2 := e.New("e2", fields.F("f2", "v2"))

		assert.NotSame(t, e1, e1.Wrap(e2))
		assert.NotSame(t, e2, e1.Wrap(e2))
		assert.NoError(t, errors.Unwrap(e1.Wrap(e2)), "errors unwrap returns nil")
		assert.Equal(t, "e1 (f1=v1): e2 (f2=v2)", e1.Wrap(e2).Error())
		assert.Equal(t, "e1 (f1=v1) (f3=v3): e2 (f2=v2)", e1.Wrap(e2, fields.F("f3", "v3")).Error())

		assert.Equal(t, "e1 (f1=v1): error(nil)", e1.Wrap(nil).Error())
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
		assert.NotErrorIs(t, e4.WithField("f1", "v1"), e4.WithField("f1", "v1")) //nolint:testifylint

		var (
			e5 error = e.NewFrom("e5", e4.WithField("f4", "v4")) // e5: e0 (f4=v4)
		)

		assert.NotErrorIs(t, e5, e4.WithField("f4", "v4"))
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

		assert.False(t, errors.As(e1, &target)) //nolint:testifylint //there is no `assert.NotErrorAs`

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
		assert.NotEqual(t, reflect.ValueOf(e1.Fields()).UnsafePointer(), reflect.ValueOf(e2.Fields()).UnsafePointer())
		assert.ElementsMatch(t, e1.Fields(), e2.Fields())
	})

	t.Run(".WithFields()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1", fields.F("f1", "v1"))
		e2 := e1.WithFields(fields.F("f2", "v2"))

		assert.NotSame(t, e1, e2)
		assert.ErrorIs(t, e2, e1, "child error passes errors.Is for parent")

		assert.NotEqual(t, reflect.ValueOf(e1.Fields()).UnsafePointer(), reflect.ValueOf(e2.Fields()).UnsafePointer())
		assert.Len(t, e1.Fields(), 1)
		assert.Len(t, e2.Fields(), 1)
		assert.Equal(t, "e1 (f1=v1) (f2=v2)", e2.Error())
	})

	t.Run(".WithField()", func(t *testing.T) {
		t.Parallel()

		e1 := e.New("e1", fields.F("f1", "v1"))
		e2 := e1.WithField("f2", "v2")

		assert.NotSame(t, e1, e2)
		assert.ErrorIs(t, e2, e1, "child error passes errors.Is for parent")

		assert.NotEqual(t, reflect.ValueOf(e1.Fields()).UnsafePointer(), reflect.ValueOf(e2.Fields()).UnsafePointer())
		assert.Len(t, e1.Fields(), 1)
		assert.Len(t, e2.Fields(), 1)
		assert.Equal(t, "e1 (f1=v1) (f2=v2)", e2.Error())
	})
}

type myErr struct {
	err string
}

func (err *myErr) Error() string {
	return err.err
}
