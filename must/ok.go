package must

import (
	"dev.gaijin.net/go/golib/e"
)

// OK is a wrapper function to simplify code, in situations where developer is sure
// that the function being called cannot return an error. E.g.
//
//	u, err := url.Parse("example.com")
//	if err != nil {
//	    panic(err)
//	}
//
// developer doesn't expect error here because the parsed path is a static and correct
//
//	u = must.OK(url.Parse("example.com"))
func OK[T any](v T, err error) T { //nolint:ireturn
	if err != nil {
		panic(e.NewFrom("no error assurance failed", err))
	}

	return v
}
