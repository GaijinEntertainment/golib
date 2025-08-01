package must

import (
	"dev.gaijin.team/go/golib/e"
)

// True returns the value if condition is true, otherwise panics.
//
//	result := must.True(logger.FromCtx(ctx))
func True[T any](v T, condition bool) T { //nolint:ireturn
	if !condition {
		panic(e.New("must.True assertion failed"))
	}

	return v
}
