// Package e provides custom error type and utilities to work with it.
//
// Its main type [Err] provides the ability to conveniently work with chains of
// errors, offering methods to wrap any errors and pass metadata using fields
// (fields.Field).
//
// Example:
//
//	var (
//		ErrJsonParseFailed = e.New("JSON parse failed")
//	)
//
//	var val any
//	if err := json.Unmarshal([]byte(`["invalid", "json]`), &val); err != nil {
//		return ErrJsonParseFailed.Wrap(err) // "JSON parse failed: unexpected end of JSON input"
//	}
//
// Any incoming error can be granted the functionality of [Err] by using the
// [From] function. Note that this is not a form of error wrapping - immediately
// unwrapping such an error will not return the original one.
//
// Example:
//
//	e.From(errors.New("error")) // "error"
//	errors.Unwrap(e.From(errors.New("error"))) // nil
//
// [Err] serializes to following format:
//
//	`<error reason>[ (fields...)][: wrapped error string]`
//
// In other words, arbitrary fields are always enclosed in parentheses, and the
// wrapped error is separated with a colon and a space.
//
// Package does not provide any methods to modify existing errors, as it is
// considered error-prone. Instead, any public method of the package that returns
// an error - returns a new instance.
//
// Although [Err] implements methods [Err.Unwrap], [Err.Is], and [Err.As], these
// methods are intended for internal use. These methods are marked as deprecated
// to avoid confusion and improve developer experience. Consider using the
// corresponding methods of the [errors] package instead.
package e
