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

func newAdapter(hook *logrusHook) *logrusadapter.Adapter {
	ll := logrus.New()
	ll.Level = logrus.TraceLevel
	ll.AddHook(hook)

	ll.SetOutput(&discardingWriter{})

	return logrusadapter.New(logrus.NewEntry(ll))
}

func TestAdapter(t *testing.T) {
	t.Parallel()

	t.Run(".Log()", func(t *testing.T) {
		t.Parallel()

		hook := logrusHook{}
		adapter := newAdapter(&hook)

		adapter.Log(42, "foo", nil)

		require.Len(t, hook, 2)
		assert.Equal(t, logrus.ErrorLevel, hook[0].Level)
		assert.Nil(t, hook[0].Data[logrus.ErrorKey])
		assert.Equal(t, logrus.InfoLevel, hook[1].Level)
		assert.Nil(t, hook[1].Data[logrus.ErrorKey])

		hook.Reset()

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

			adapter.Log(tc.level, tc.msg, tc.err)

			require.Len(t, hook, 1)
			assert.Equal(t, tc.logrusLevel, hook[0].Level)
			assert.Equal(t, tc.msg, hook[0].Message)
			assert.Same(t, tc.err, hook[0].Data[logrus.ErrorKey])
		}
	})

	t.Run(".WithFields()", func(t *testing.T) {
		t.Parallel()

		hook := logrusHook{}
		adapter := newAdapter(&hook).WithFields(fields.F("foo", "bar"), fields.F("baz", 42))

		adapter.Log(logger.LevelInfo, "test", nil)
		adapter.Log(logger.LevelDebug, "test", nil, fields.F("bux", "qux"))
		adapter.Log(logger.LevelDebug, "test", nil, fields.F("baz", "bar"))

		require.Len(t, hook, 3)

		assert.Equal(t, "bar", hook[0].Data["foo"])
		assert.Equal(t, 42, hook[0].Data["baz"])

		assert.Equal(t, "bar", hook[1].Data["foo"])
		assert.Equal(t, 42, hook[1].Data["baz"])
		assert.Equal(t, "qux", hook[1].Data["bux"])

		assert.Equal(t, "bar", hook[2].Data["foo"])
		assert.Equal(t, "bar", hook[2].Data["baz"])
	})

	t.Run(".WithName()", func(t *testing.T) {
		t.Parallel()

		hook := logrusHook{}
		adapter := newAdapter(&hook).WithName("test-logger")

		adapter.Log(logger.LevelInfo, "test", nil)

		require.Len(t, hook, 1)

		assert.Equal(t, "test-logger", hook[0].Data[logrusadapter.LoggerNameKey])
	})
}
