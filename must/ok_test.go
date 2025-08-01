package must_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/must"
)

func TestOK(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		uri := must.OK(url.Parse("https://example.com"))
		assert.Equal(t, "https://example.com", uri.String())
	})

	imminentError := func() (bool, error) { return false, e.New("imminent error") }

	err := catchPanicError(t, func() {
		_ = must.OK(imminentError())
	})

	assert.ErrorContains(t, err, "must.OK assertion failed: imminent error")
}

func catchPanicError(t *testing.T, fn func()) (err error) {
	t.Helper()

	defer func() {
		if rv := recover(); rv != nil {
			re, ok := rv.(error)
			if !ok {
				t.Fatal("recovered panic value is not an error")
			}

			err = re
		}
	}()

	fn()

	return nil
}
