package zapadapter_test

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type discardingWriter struct{}

func (*discardingWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (*discardingWriter) Sync() (err error) {
	return nil
}

func zapLogger() (*zap.Logger, *observer.ObservedLogs) {
	discardingCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		&discardingWriter{},
		zapcore.DebugLevel,
	)

	testCore, logs := observer.New(zapcore.DebugLevel)

	return zap.New(zapcore.NewTee(discardingCore, testCore)), logs
}

func TestAdapter(t *testing.T) {
	t.Parallel()

	_, _ = zapLogger()
}
