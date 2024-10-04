package logger

import (
	"context"
)

type CtxKey struct{}

// FromCtx attempts to extract logger from the context. In case logger is not
// present, or it is not of the correct type, it returns empty logger and false.
func FromCtx(ctx context.Context) (Logger, bool) {
	logger, ok := ctx.Value(CtxKey{}).(Logger)

	return logger, ok
}

// FromCtxOrNop extracts logger from the context. In case logger is not
// present, or it is not of the correct type, it returns a nop logger.
func FromCtxOrNop(ctx context.Context) Logger {
	if logger, ok := FromCtx(ctx); ok {
		return logger
	}

	return NewNop()
}

// ToCtx stores logger in the context, it can be extracted with [FromCtx] or [FromCtxOrNop].
func ToCtx(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, CtxKey{}, logger)
}
