package bufferadapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger/bufferadapter"
)

func TestAdapter(t *testing.T) {
	t.Parallel()

	t.Run(".Log()", func(t *testing.T) {
		t.Parallel()

		buff := bufferadapter.LogEntries{}
		adapter := bufferadapter.New(&buff)

		adapter.Log(42, "foo", nil)
		require.Len(t, buff, 1)
		assert.Equal(t, bufferadapter.LogEntry{
			Level: 42,
			Msg:   "foo",
		}, buff[0])

		err := e.New("some error")
		adapter.Log(42, "foo", err, fields.F("foo", "bar"))
		require.Len(t, buff, 2)
		assert.Equal(t, bufferadapter.LogEntry{
			Level:  42,
			Msg:    "foo",
			Error:  err,
			Fields: fields.List{fields.F("foo", "bar")},
		}, buff[1])

		buff.Reset()
		assert.Empty(t, buff)
	})

	t.Run(".WithFields()", func(t *testing.T) {
		t.Parallel()

		buff := bufferadapter.LogEntries{}
		adapterSrc := bufferadapter.New(&buff)
		adapter := adapterSrc.WithName("my-logger").WithFields(fields.F("foo", "bar"))

		require.NotSame(t, adapterSrc, adapter)

		adapter.Log(42, "foo", nil)
		require.Len(t, buff, 1)
		assert.Equal(t, bufferadapter.LogEntry{
			LoggerName: "my-logger",
			Level:      42,
			Msg:        "foo",
			Fields:     fields.List{fields.F("foo", "bar")},
		}, buff[0])

		adapter.Log(42, "foo", nil, fields.F("baz", "qux"))
		require.Len(t, buff, 2)
		assert.Equal(t, bufferadapter.LogEntry{
			LoggerName: "my-logger",
			Level:      42,
			Msg:        "foo",
			Fields:     fields.List{fields.F("foo", "bar"), fields.F("baz", "qux")},
		}, buff[1])
	})

	t.Run(".WithName()", func(t *testing.T) {
		t.Parallel()

		buff := bufferadapter.LogEntries{}
		adapterSrc := bufferadapter.New(&buff)
		adapter := adapterSrc.WithFields(fields.F("foo", "bar")).WithName("my-logger")

		require.NotSame(t, adapterSrc, adapter)

		adapter.Log(42, "foo", nil)
		require.Len(t, buff, 1)
		assert.Equal(t, bufferadapter.LogEntry{
			LoggerName: "my-logger",
			Level:      42,
			Msg:        "foo",
			Fields:     fields.List{fields.F("foo", "bar")},
		}, buff[0])
	})
}
