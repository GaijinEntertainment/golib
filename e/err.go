package e

import (
	"errors"
	"fmt"
	"strings"

	"github.com/GaijinEntertainment/golib/fields"
)

type Err struct {
	err     error
	wrapped error
	fields  fields.List
}

// New creates new instance of Err.
func New(reason string, f ...fields.Field) *Err {
	return &Err{
		err:     errors.New(reason), //nolint:err113
		wrapped: nil,
		fields:  f,
	}
}

// NewFrom creates new instance of Err that wraps origin error.
func NewFrom(origin error, reason string, f ...fields.Field) *Err {
	return &Err{
		err:     errors.New(reason), //nolint:err113
		wrapped: origin,
		fields:  f,
	}
}

// From transforms existing error to Err.
func From(origin error) *Err {
	return &Err{
		err:     origin,
		wrapped: nil,
		fields:  nil,
	}
}

func (e *Err) Error() string {
	b := &strings.Builder{}
	writeTo(b, e)

	return b.String()
}

func writeTo(b *strings.Builder, err error) {
	if b.Len() > 0 {
		b.WriteString(": ")
	}

	ee, ok := err.(*Err) //nolint:errorlint
	if !ok {
		b.WriteString(err.Error())
		return
	}

	if ee == nil || ee.err == nil {
		b.WriteString("nil")
		return
	}

	b.WriteString(ee.err.Error())

	if ee.fields != nil {
		b.WriteRune(' ')
		ee.fields.WriteTo(b)
	}

	if ee.wrapped != nil {
		writeTo(b, ee.wrapped)
	}
}

// Wrap wraps provided errors into a single Err in order. Strings will be
// constructed to errors using [New]. In case of non-error and non-string values,
// it will be converted to a string using fmt.Sprintf("%#v").
//
// Example:
//
//	e.Wrap(e.New("e1", fields.F("f1", "v1")), errors.New("e2"), "e3") // e1 (f1=v1): e2: e3
func Wrap(args ...any) *Err {
	var err *Err

	for i := len(args) - 1; i >= 0; i-- {
		var er *Err

		switch v := args[i].(type) {
		case error:
			er = From(v)

		case string:
			er = New(v)

		default:
			er = New(fmt.Sprintf("%T(%v)", v, v))
		}

		if err != nil {
			er.wrapped = err
		}

		err = er
	}

	return err
}

// Wrap creates new instance of Err that wraps provided error with source one.
func (e *Err) Wrap(err error, f ...fields.Field) *Err {
	return &Err{
		err:     e,
		wrapped: err,
		fields:  f,
	}
}

// Unwrap returns wrapped error. Returns nil in case there is no wrapper error.
func (e *Err) Unwrap() error {
	return e.wrapped
}

// Is reports whether any error in the chain matches the target error.
//
// This method implemented only to satisfy errors.Is interface, for checking
// errors use [errors.Is] instead.
func (e *Err) Is(tgt error) bool {
	return errors.Is(e.err, tgt) || errors.Is(e.wrapped, tgt)
}

// As finds the first error in err's tree that matches target, and if one is
// found, sets target to that error value and returns true. Otherwise, it returns
// false.
//
// This method implemented only to satisfy errors.Is interface, for checking
// errors use [errors.Is] instead.
func (e *Err) As(target any) bool {
	if e.err != nil && errors.As(e.err, target) {
		return true
	}

	if e.wrapped != nil && errors.As(e.wrapped, target) {
		return true
	}

	return false
}
