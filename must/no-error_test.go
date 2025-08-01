package must_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/must"
)

func TestNoErr(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		var data []string

		must.NoErr(json.Unmarshal([]byte(`["valid", "json"]`), &data))
		assert.Equal(t, []string{"valid", "json"}, data)
	})

	imminentError := func() error { return e.New("imminent error") }

	err := catchPanicError(t, func() {
		must.NoErr(imminentError())
	})

	assert.ErrorContains(t, err, "must.NoErr assertion failed: imminent error")
}
