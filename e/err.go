package e

import (
	"errors"
	"slices"
	"strings"

	"dev.gaijin.team/go/golib/fields"
)

type Err struct {
	reason          string
	origin, wrapped error
	fields          fields.List
}

// New creates new instance of Err.
func New(reason string, f ...fields.Field) *Err {
	return &Err{
		reason:  reason,
		origin:  nil,
		wrapped: nil,
		fields:  f,
	}
}

// NewFrom creates new instance of Err that wraps origin error.
func NewFrom(reason string, origin error, f ...fields.Field) *Err {
	return &Err{
		reason:  reason,
		origin:  nil,
		wrapped: origin,
		fields:  f,
	}
}

// From transforms existing error to Err. It is not wrapping operation - unwrap
// will not return origin error. Passing nil to this function will result with
// empty error.
func From(origin error, f ...fields.Field) *Err {
	return &Err{
		reason:  "",
		origin:  origin,
		wrapped: nil,
		fields:  f,
	}
}

// Error returns string representation of the error.
func (e *Err) Error() string {
	b := &strings.Builder{}
	writeTo(b, e)

	return b.String()
}

func writeTo(b *strings.Builder, err error) {
	ee, ok := err.(*Err) //nolint:errorlint
	if !ok {
		if b.Len() > 0 {
			b.WriteString(": ")
		}

		b.WriteString(err.Error())

		return
	}

	if str := errString(ee); str != "" {
		if b.Len() > 0 {
			b.WriteString(": ")
		}

		b.WriteString(str)
	}

	if ee == nil {
		return
	}

	if ee.origin != nil {
		writeTo(b, ee.origin)
	}

	if len(ee.fields) > 0 {
		b.WriteRune(' ')
		ee.fields.WriteTo(b)
	}

	if ee.wrapped != nil {
		writeTo(b, ee.wrapped)
	}
}

func errString(e *Err) string {
	if e == nil {
		return "(*e.Err)(nil)"
	}

	if e.reason != "" {
		return e.reason
	}

	if e.origin == nil {
		return "(*e.Err)(empty)"
	}

	return ""
}

// Wrap provided errors into each other in order, resulting with singular error.
// Passed errors will be converted to [Err] using [From] function.
//
// If no errors provided, (*e.Err)(nil) will be returned.
//
// Example:
//
//	e.Wrap(errors.New("e1"), errors.New("e2"), errors.New("e3")) // e1: e2: e3
func Wrap(args ...error) (err *Err) {
	for i := len(args) - 1; i >= 0; i-- {
		er := From(args[i])

		if err != nil {
			er.wrapped = err
		}

		err = er
	}

	return err
}

// Clone creates a new instance of Err with the same error, wrapped error, and
// cloned fields container.
func (e *Err) Clone() *Err {
	return &Err{
		reason:  e.reason,
		origin:  e.origin,
		wrapped: e.wrapped,
		fields:  slices.Clone(e.fields),
	}
}

// WithFields creates a clone of initial error with fields added to it. Fields of
// initial error will not be changed.
func (e *Err) WithFields(f ...fields.Field) *Err {
	ee := e.Clone()

	ee.fields.Add(f...)

	return ee
}

// WithField alike [Err.WithFields], but creates an error with single field added.
func (e *Err) WithField(key string, val any) *Err {
	return e.WithFields(fields.F(key, val))
}

// Fields returns fields of the error.
func (e *Err) Fields() fields.List {
	return e.fields
}

// Reason returns reason string of the error without fields and wrapped errors.
func (e *Err) Reason() string {
	if e.reason == "" {
		return e.origin.Error()
	}

	return e.reason
}

// Wrap creates new instance of Err that wraps provided error with source one.
//
// Example:
//
//	e.New("e1").Wrap(errors.New("e2")) // e1: e2
func (e *Err) Wrap(err error, f ...fields.Field) *Err {
	return &Err{
		reason:  "",
		origin:  e,
		wrapped: err,
		fields:  f,
	}
}

// Unwrap implemented only for purposes of compatibility with [errors.Unwrap].
//
// Deprecated: use [errors.Unwrap] instead.
func (e *Err) Unwrap() error {
	return e.wrapped
}

// Is implemented only for purposes of compatibility with [errors.Is]. It only
// checks error's origin. Check of wrapped errors performed by [errors.Is]
// itself.
//
// Deprecated: use [errors.Is] instead.
func (e *Err) Is(err error) bool {
	return errors.Is(e.origin, err)
}

// As implemented only for purposes of compatibility with [errors.As]. It only
// checks error's origin. Check of wrapped errors performed by [errors.As]
// itself.
//
// Deprecated: use [errors.As] instead.
func (e *Err) As(target any) bool {
	return errors.As(e.origin, target)
}
