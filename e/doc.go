// Package e provides custom error type and utilities to work with it.
//
// It's main type [Err] provides ability to conveniently wrap any errors along
// with arbitrary amount of fields (fields.Field).
//
// Example:
//
//	var (
//		ErrJsonParseFailed = e.New("JSON parse failed")
//	)
//
//	var val any
//	if err := json.Unmarshal([]byte(`["invalid", "json]`), &val); err != nil {
//		return ErrJsonParseFailed.Wrap(err) // JSON parse failed: unexpected end of JSON input
//	}
//
// It serializes error to following format:
//
//	`<error reason>[ (fields...)][: wrapped error string]`
//
// In other words arbitrary fields are always in parentheses, wrapped error is
// separated from error's reason with colon.
//
// [Wrap] function that casts provided values to errors and wraps them into each
// other consecutively.
//
// Example:
//
//	err := e.Wrap("error", e.New("wrapped error"), 42) // error: wrapped error: int(42)
//
// Package does not provide any methods to modify existing errors, as it is seen
// as error-prone. Instead, each methods returns new instance of [Err].
package e
