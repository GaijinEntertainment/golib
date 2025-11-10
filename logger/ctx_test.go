package logger_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/logger"
	"dev.gaijin.team/go/golib/logger/bufferadapter"
)

func TestToAndFromCtx(t *testing.T) {
	t.Parallel()

	adapter, _ := bufferadapter.New()
	lgr := logger.New(adapter)

	ctxEmpty := context.Background()
	ctxLogger := logger.ToCtx(ctxEmpty, lgr)

	require.NotEqual(t, ctxLogger, ctxEmpty, "ToCtx creates a new context")

	{
		lgrCtx, ok := logger.FromCtx(ctxEmpty)
		assert.False(t, ok, "should return false for empty context")
		assert.True(t, lgrCtx.IsZero(), "should return zero-value logger for empty context")
	}

	{
		lgrCtx, ok := logger.FromCtx(ctxLogger)
		assert.True(t, ok, "should return true for context with logger")
		assert.False(t, lgrCtx.IsZero(), "should not return zero-value logger for context with logger")
		assert.True(t, logger.IsEqual(lgr, lgrCtx), "should return same logger as stored in context")
	}
}

func TestFromCtxOrNop(t *testing.T) {
	t.Parallel()

	adapter, _ := bufferadapter.New()
	lgr := logger.New(adapter)

	ctxEmpty := context.Background()
	ctxLogger := logger.ToCtx(ctxEmpty, lgr)

	{
		lgrCtx := logger.FromCtxOrNop(ctxEmpty)
		assert.True(t, lgrCtx.IsNop(), "should return nop logger for empty context")
	}

	{
		lgrCtx := logger.FromCtxOrNop(ctxLogger)
		assert.False(t, lgrCtx.IsNop(), "should not return nop logger for context with logger")
		assert.True(t, logger.IsEqual(lgr, lgrCtx), "should return same logger as stored in context")
	}
}
