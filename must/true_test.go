package must_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.gaijin.team/go/golib/logger"
	"dev.gaijin.team/go/golib/must"
)

func TestTrue(t *testing.T) {
	t.Parallel()

	assert.NotPanics(t, func() {
		lgr := logger.NewNop()

		result := must.True(logger.FromCtx(logger.ToCtx(t.Context(), lgr)))
		assert.Equal(t, lgr, result)
	})

	err := catchPanicError(t, func() {
		_ = must.True(logger.FromCtx(t.Context()))
	})

	assert.ErrorContains(t, err, "must.True assertion failed")
}
