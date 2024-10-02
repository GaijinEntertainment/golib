package must

import (
	"github.com/GaijinEntertainment/golib/e"
)

// NoErr is a wrapper function to simplify code, in situations where developer is sure
// that the function being called cannot return an error. E.g.
//
//	err := json.Unmarshal([]byte(`["totally", "valid", "data"]`), &data)
//	if err != nil {
//	    panic(err)
//	}
//
// developer doesn't expect error here because the parsed data is static and correct
//
//	must.OK(json.Unmarshal([]byte(`["totally", "valid", "data"]`), &data))
func NoErr(err error) {
	if err != nil {
		panic(e.NewFrom("OK assurance failed", err))
	}
}
