// Package must provides assertion functions that panic on failure.
//
// # When to Use
//
// Use must functions when you are certain that operations should not fail
// and want to eliminate error handling boilerplate in situations where
// errors indicate programming bugs rather than runtime conditions.
//
// Common scenarios:
//   - Parsing static/hardcoded data (URLs, JSON, regex patterns)
//   - Operations that should always succeed in your specific context
//   - Configuration parsing where failure means misconfiguration
//
// # When NOT to Use
//
// Avoid must functions for:
//   - User input validation
//   - Network operations
//   - File I/O operations
//   - Any operation that can legitimately fail at runtime
//
// # Functions
//
//   - must.OK[T](value T, err error) T - returns value if err is nil, panics otherwise
//   - must.True[T](value T, condition bool) T - returns value if condition is true, panics otherwise
//   - must.NoErr(err error) - panics if err is not nil
//
// # Examples
//
//	// Static data parsing
//	baseURL := must.OK(url.Parse("https://api.example.com"))
//	pattern := must.OK(regexp.Compile(`^[a-z]+$`))
//
//	// Configuration validation
//	config := must.True(loadConfig(), config.IsValid())
//
//	// JSON unmarshaling of known-good data
//	var data Config
//	must.NoErr(json.Unmarshal(embeddedConfig, &data))
package must
