package bufferadapter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger/bufferadapter"
)

func Test_New(t *testing.T) {
	t.Parallel()

	adapter, buff := bufferadapter.New()
	require.NotNil(t, adapter)
	require.NotNil(t, buff)

	adapter.Log(42, "test", nil)
	require.Equal(t, 1, buff.Len())
	assert.Equal(t, "test", buff.Get(0).Msg)
}

func TestAdapter(t *testing.T) {
	t.Parallel()

	t.Run(".Log()", func(t *testing.T) {
		t.Parallel()

		adapter, buff := bufferadapter.New()

		adapter.Log(42, "foo", nil)
		require.Equal(t, 1, buff.Len())
		assert.Equal(t, bufferadapter.LogEntry{
			Level: 42,
			Msg:   "foo",
		}, buff.Get(0))

		err := e.New("some error")
		adapter.Log(42, "foo", err, fields.F("foo", "bar"))
		require.Equal(t, 2, buff.Len())
		assert.Equal(t, bufferadapter.LogEntry{
			Level:  42,
			Msg:    "foo",
			Error:  err,
			Fields: fields.List{fields.F("foo", "bar")},
		}, buff.Get(1))

		buff.Reset()
		assert.Equal(t, 0, buff.Len())
	})

	t.Run(".WithFields()", func(t *testing.T) {
		t.Parallel()

		adapterSrc, buff := bufferadapter.New()
		adapter := adapterSrc.WithName("my-logger").WithFields(fields.F("foo", "bar"))

		require.NotSame(t, adapterSrc, adapter)

		adapter.Log(42, "foo", nil)
		require.Equal(t, 1, buff.Len())
		assert.Equal(t, bufferadapter.LogEntry{
			LoggerName: "my-logger",
			Level:      42,
			Msg:        "foo",
			Fields:     fields.List{fields.F("foo", "bar")},
		}, buff.Get(0))

		adapter.Log(42, "foo", nil, fields.F("baz", "qux"))
		require.Equal(t, 2, buff.Len())
		assert.Equal(t, bufferadapter.LogEntry{
			LoggerName: "my-logger",
			Level:      42,
			Msg:        "foo",
			Fields:     fields.List{fields.F("foo", "bar"), fields.F("baz", "qux")},
		}, buff.Get(1))
	})

	t.Run(".WithName()", func(t *testing.T) {
		t.Parallel()

		adapterSrc, buff := bufferadapter.New()
		adapter := adapterSrc.WithFields(fields.F("foo", "bar")).WithName("my-logger")

		require.NotSame(t, adapterSrc, adapter)

		adapter.Log(42, "foo", nil)
		require.Equal(t, 1, buff.Len())
		assert.Equal(t, bufferadapter.LogEntry{
			LoggerName: "my-logger",
			Level:      42,
			Msg:        "foo",
			Fields:     fields.List{fields.F("foo", "bar")},
		}, buff.Get(0))
	})
}

func TestLogEntries(t *testing.T) {
	t.Parallel()

	t.Run("GetAll()", func(t *testing.T) {
		t.Parallel()

		adapter, buff := bufferadapter.New()

		adapter.Log(10, "first", nil)
		adapter.Log(20, "second", nil)
		adapter.Log(30, "third", nil)

		entries := buff.GetAll()
		require.Len(t, entries, 3)
		assert.Equal(t, "first", entries[0].Msg)
		assert.Equal(t, "second", entries[1].Msg)
		assert.Equal(t, "third", entries[2].Msg)
	})

	t.Run("Get() panics on out of range", func(t *testing.T) {
		t.Parallel()

		_, buff := bufferadapter.New()

		assert.Panics(t, func() {
			buff.Get(0)
		})
	})

	t.Run("Reset() empties the buffer", func(t *testing.T) {
		t.Parallel()

		adapter, buff := bufferadapter.New()

		for i := range 10 {
			adapter.Log(i, "test", nil)
		}

		require.Equal(t, 10, buff.Len())

		buff.Reset()
		assert.Equal(t, 0, buff.Len())
	})
}
