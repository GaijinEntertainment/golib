package e

import (
	"dev.gaijin.team/go/golib/fields"
)

type ErrorLogger func(msg string, err error, fs ...fields.Field)

// Log logs the provided error using the given logger.
//
// If the error is nil, the function does nothing.
//
// If the error is of type [Err], its reason is used as the error message, the
// wrapped error is passed as the actual error, and the error's fields are passed
// as log fields.
//
// Otherwise, err.Error() is used as the error message, and nil is passed as the
// actual error.
func Log(err error, f ErrorLogger) {
	if err == nil {
		return
	}

	// we're not interested in wrapped error, therefore we're only typecasting it.
	if e, ok := err.(*Err); ok { //nolint:errorlint
		var wrapped error
		if len(e.errs) > 1 {
			wrapped = e.errs[1]
		}

		f(e.Reason(), wrapped, e.fields...)

		return
	}

	f(err.Error(), nil)
}
