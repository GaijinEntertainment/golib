package must_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.gaijin.team/go/golib/must"
)

func TestOK(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		_ = must.OK(regexp.Compile("[a-z]")) //nolint:gocritic
	})

	assert.Panics(t, func() {
		_ = must.OK(regexp.Compile("[a-z")) //nolint:staticcheck,gocritic
	})

	err := catchPanic(t, func() {
		_ = must.OK(errorFn())
	})

	assert.ErrorIs(t, err, errTestError)
	assert.Equal(t, "no error assurance failed: test error", err.Error())
}

var errTestError = errors.New("test error")

func errorFn() (bool, error) {
	return false, errTestError
}

func catchPanic(t *testing.T, panickingFn func()) (err error) {
	t.Helper()

	defer func() {
		err = recover().(error) //nolint:forcetypeassert
	}()

	panickingFn()

	return nil
}
