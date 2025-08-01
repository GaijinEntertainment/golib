package must

import (
	"dev.gaijin.team/go/golib/e"
)

// OK returns the value if error is nil, otherwise panics.
//
//	u := must.OK(url.Parse("https://example.com"))
func OK[T any](v T, err error) T { //nolint:ireturn
	if err != nil {
		panic(e.NewFrom("must.OK assertion failed", err))
	}

	return v
}
