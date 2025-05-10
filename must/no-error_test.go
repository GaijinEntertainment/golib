package must_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.gaijin.team/go/golib/must"
)

func TestNoErr(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		var v []string

		must.NoErr(json.Unmarshal([]byte(`["foo"]`), &v))
	})

	assert.Panics(t, func() {
		var v []string

		must.NoErr(json.Unmarshal([]byte(`["foo]`), &v))
	})

	err := catchPanic(t, func() {
		must.NoErr(errorFn1())
	})

	assert.ErrorIs(t, err, errTestError)
	assert.Equal(t, "NoErr assurance failed: test error", err.Error())
}

func errorFn1() error {
	return errTestError
}
