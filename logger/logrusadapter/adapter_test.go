//nolint:err113
package logrusadapter_test

import (
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
	"dev.gaijin.team/go/golib/logger/logrusadapter"
)

type discardingWriter struct{}

func (*discardingWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type logrusHook []*logrus.Entry

func (*logrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *logrusHook) Fire(entry *logrus.Entry) error {
	*h = append(*h, entry)

	return nil
}

func (h *logrusHook) Reset() {
	*h = logrusHook{}
}

func newAdapter(hook *logrusHook, opts ...logrusadapter.Option) *logrusadapter.Adapter {
	ll := logrus.New()

	ll.Level = logrus.TraceLevel
	ll.AddHook(hook)
	ll.SetOutput(&discardingWriter{})

	return logrusadapter.New(logrus.NewEntry(ll), opts...)
}

func TestAdapter(t *testing.T) {
	t.Parallel()

	t.Run(".Log() with standard levels", func(t *testing.T) {
		t.Parallel()

		hook := logrusHook{}
		adapter := newAdapter(&hook)

		tt := []struct {
			level       int
			logrusLevel logrus.Level
			msg         string
			err         error
		}{
			{logger.LevelDebug, logrus.DebugLevel, "debug", errors.New("debug")},
			{logger.LevelInfo, logrus.InfoLevel, "info", errors.New("info")},
			{logger.LevelWarning, logrus.WarnLevel, "warning", errors.New("warning")},
			{logger.LevelError, logrus.ErrorLevel, "error", errors.New("error")},
			{logger.LevelTrace, logrus.TraceLevel, "trace", errors.New("trace")},
		}

		for _, tc := range tt {
			hook.Reset()

			adapter.Log(tc.level, tc.msg, fields.F("error", tc.err))

			require.Len(t, hook, 1)
			assert.Equal(t, tc.logrusLevel, hook[0].Level)
			assert.Equal(t, tc.msg, hook[0].Message)
			assert.Same(t, tc.err, hook[0].Data["error"])
		}
	})

	t.Run(".Log() with unknown level", func(t *testing.T) {
		t.Parallel()

		hook := logrusHook{}
		adapter := newAdapter(&hook)

		adapter.Log(42, "unknown level")

		// Should map to InfoLevel by default, creating only ONE entry
		require.Len(t, hook, 1)
		assert.Equal(t, logrus.InfoLevel, hook[0].Level)
		assert.Equal(t, "unknown level", hook[0].Message)
		assert.Nil(t, hook[0].Data[logrus.ErrorKey])
	})

	t.Run(".WithFields()", func(t *testing.T) {
		t.Parallel()

		hook := logrusHook{}
		adapter := newAdapter(&hook).WithFields(fields.F("foo", "bar"), fields.F("baz", 42))

		adapter.Log(logger.LevelInfo, "test")
		adapter.Log(logger.LevelDebug, "test", fields.F("bux", "qux"))
		adapter.Log(logger.LevelDebug, "test", fields.F("baz", "bar"))

		require.Len(t, hook, 3)

		assert.Equal(t, "bar", hook[0].Data["foo"])
		assert.Equal(t, 42, hook[0].Data["baz"])

		assert.Equal(t, "bar", hook[1].Data["foo"])
		assert.Equal(t, 42, hook[1].Data["baz"])
		assert.Equal(t, "qux", hook[1].Data["bux"])

		assert.Equal(t, "bar", hook[2].Data["foo"])
		assert.Equal(t, "bar", hook[2].Data["baz"])
	})

	t.Run("WithLogLevelMapper option", func(t *testing.T) {
		t.Parallel()

		hook := logrusHook{}

		// Custom mapper that always returns ErrorLevel
		customMapper := func(level int) logrus.Level {
			return logrus.ErrorLevel
		}

		adapter := newAdapter(&hook, logrusadapter.WithLogLevelMapper(customMapper))

		// Log at Info level, but should appear as Error due to custom mapper
		adapter.Log(logger.LevelInfo, "test")

		require.Len(t, hook, 1)
		assert.Equal(t, logrus.ErrorLevel, hook[0].Level)
		assert.Equal(t, "test", hook[0].Message)
	})

	t.Run("custom mapper preserved in derived adapters", func(t *testing.T) {
		t.Parallel()

		hook := logrusHook{}

		// Custom mapper that maps everything to WarnLevel
		customMapper := func(level int) logrus.Level {
			return logrus.WarnLevel
		}

		adapter := newAdapter(&hook, logrusadapter.WithLogLevelMapper(customMapper))
		derived := adapter.WithFields(fields.F("foo", "bar"))

		// Both should use custom mapper
		adapter.Log(logger.LevelInfo, "original")
		derived.Log(logger.LevelDebug, "derived")

		require.Len(t, hook, 2)
		assert.Equal(t, logrus.WarnLevel, hook[0].Level)
		assert.Equal(t, logrus.WarnLevel, hook[1].Level)
		assert.Equal(t, "bar", hook[1].Data["foo"])
	})
}
