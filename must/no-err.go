package must

import (
	"dev.gaijin.team/go/golib/e"
)

// NoErr panics if error is not nil.
//
//	must.NoErr(json.Unmarshal(staticJSON, &data))
func NoErr(err error) {
	if err != nil {
		panic(e.NewFrom("must.NoErr assertion failed", err))
	}
}
