package e

import (
	"errors"
	"slices"
	"strings"

	"dev.gaijin.team/go/golib/fields"
)

type Err struct {
	errs   []error
	fields fields.List
}

// New creates new instance of Err.
func New(reason string, f ...fields.Field) *Err {
	return From(errors.New(reason), f...) //nolint:err113
}

// NewFrom creates new instance of Err that wraps origin error.
func NewFrom(reason string, wrapped error, f ...fields.Field) *Err {
	if wrapped == nil {
		return New(reason, f...)
	}

	return &Err{
		errs:   []error{errors.New(reason), wrapped}, //nolint:err113
		fields: f,
	}
}

// From transforms existing error to Err. It is not wrapping operation - unwrap
// will not return origin error. Passing nil to this function will result with
// empty error.
func From(origin error, f ...fields.Field) *Err {
	if origin == nil {
		origin = errors.New("error(nil)") //nolint:err113
	}

	return &Err{
		errs:   []error{origin},
		fields: f,
	}
}

// Wrap creates new instance of Err that wraps provided error with source one.
//
// Example:
//
//	e.New("e1").Wrap(errors.New("e2")) // e1: e2
func (e *Err) Wrap(err error, f ...fields.Field) *Err {
	if err == nil {
		err = errors.New("error(nil)") //nolint:err113
	}

	return &Err{
		errs:   []error{e, err},
		fields: f,
	}
}

// Error returns string representation of the error.
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

	b.WriteString(ee.Reason())

	if ee == nil {
		return
	}

	if len(ee.fields) > 0 {
		b.WriteRune(' ')
		ee.fields.WriteTo(b)
	}

	if len(ee.errs) > 1 {
		writeTo(b, ee.errs[1])
	}
}

// Clone creates a new instance of Err with the same error, wrapped error, and
// cloned fields container.
func (e *Err) Clone() *Err {
	return &Err{
		errs:   slices.Clone(e.errs),
		fields: slices.Clone(e.fields),
	}
}

// WithFields creates new error with source error as origin and provided fields as its fields.
func (e *Err) WithFields(f ...fields.Field) *Err {
	return From(e, f...)
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
	if e == nil {
		return "(*e.Err)(nil)"
	}

	if len(e.errs) == 0 {
		return "(*e.Err)(empty)"
	}

	return e.errs[0].Error()
}

// Unwrap implemented only for purposes of compatibility with [errors.Is] and
// [errors.As] methods.
func (e *Err) Unwrap() []error {
	return e.errs
}
