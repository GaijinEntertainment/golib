//nolint:err113
package slogadapter_test

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
	"dev.gaijin.team/go/golib/logger/slogadapter"
)

type bufferEntry struct {
	Message string
	Level   slog.Level
	Attrs   map[string]any
	Group   string
}

type entriesBuffer []bufferEntry

func (b *entriesBuffer) Reset() {
	*b = []bufferEntry{}
}

type bufferHandler struct {
	buf   *entriesBuffer
	attrs []slog.Attr
	group string
}

func (*bufferHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *bufferHandler) Handle(_ context.Context, r slog.Record) error {
	entry := bufferEntry{
		Message: r.Message,
		Level:   r.Level,
		Attrs:   map[string]any{},
		Group:   h.group,
	}

	for _, a := range h.attrs {
		entry.Attrs[a.Key] = a.Value.Any()
	}

	r.Attrs(func(a slog.Attr) bool {
		entry.Attrs[a.Key] = a.Value.Any()

		return true
	})

	*h.buf = append(*h.buf, entry)

	return nil
}

func (h *bufferHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &bufferHandler{
		buf:   h.buf,
		attrs: append(slices.Clone(h.attrs), attrs...),
		group: h.group,
	}
}

func (h *bufferHandler) WithGroup(name string) slog.Handler {
	return &bufferHandler{
		buf:   h.buf,
		attrs: h.attrs,
		group: name,
	}
}

func newAdapter(buf *entriesBuffer, opts ...slogadapter.Option) *slogadapter.Adapter {
	return slogadapter.New(slog.New(&bufferHandler{buf: buf}), opts...) //nolint:exhaustruct
}

const errorKey = "error"

func TestAdapter(t *testing.T) {
	t.Parallel()

	t.Run(".Log() with standard levels", func(t *testing.T) {
		t.Parallel()

		buf := entriesBuffer{}
		adapter := newAdapter(&buf)

		tt := []struct {
			level     int
			slogLevel slog.Level
			msg       string
			err       error
		}{
			{logger.LevelDebug, slog.LevelDebug, "debug", errors.New("debug")},
			{logger.LevelInfo, slog.LevelInfo, "info", errors.New("info")},
			{logger.LevelWarning, slog.LevelWarn, "warning", errors.New("warning")},
			{logger.LevelError, slog.LevelError, "error", errors.New("error")},
			{logger.LevelTrace, slog.LevelDebug, "trace", errors.New("trace")},
		}

		for _, tc := range tt {
			buf.Reset()

			adapter.Log(tc.level, tc.msg, fields.F("error", tc.err))

			require.Len(t, buf, 1)
			assert.Equal(t, tc.slogLevel, buf[0].Level)
			assert.Equal(t, tc.msg, buf[0].Message)
			assert.Same(t, tc.err, buf[0].Attrs["error"])
		}
	})

	t.Run(".Log() with unknown level", func(t *testing.T) {
		t.Parallel()

		buf := entriesBuffer{}
		adapter := newAdapter(&buf)

		adapter.Log(42, "unknown level")

		// Should map to InfoLevel by default, creating only ONE entry
		require.Len(t, buf, 1)
		assert.Equal(t, slog.LevelInfo, buf[0].Level)
		assert.Equal(t, "unknown level", buf[0].Message)
		assert.Nil(t, buf[0].Attrs[errorKey])
	})

	t.Run(".WithFields()", func(t *testing.T) {
		t.Parallel()

		buf := entriesBuffer{}
		adapter := newAdapter(&buf).WithFields(fields.F("foo", "bar"), fields.F("baz", 42))

		adapter.Log(logger.LevelInfo, "test")
		adapter.Log(logger.LevelDebug, "test", fields.F("bux", "qux"))
		adapter.Log(logger.LevelDebug, "test", fields.F("baz", "bar"))

		require.Len(t, buf, 3)

		assert.Equal(t, "bar", buf[0].Attrs["foo"])
		assert.Equal(t, int64(42), buf[0].Attrs["baz"])

		assert.Equal(t, "bar", buf[1].Attrs["foo"])
		assert.Equal(t, int64(42), buf[1].Attrs["baz"])
		assert.Equal(t, "qux", buf[1].Attrs["bux"])

		assert.Equal(t, "bar", buf[2].Attrs["foo"])
		assert.Equal(t, "bar", buf[2].Attrs["baz"])
	})

	t.Run("WithLogLevelMapper option", func(t *testing.T) {
		t.Parallel()

		buf := entriesBuffer{}

		// Custom mapper that always returns ErrorLevel
		customMapper := func(level int) slog.Level {
			return slog.LevelError
		}

		adapter := newAdapter(&buf, slogadapter.WithLogLevelMapper(customMapper))

		// Log at Info level, but should appear as Error due to custom mapper
		adapter.Log(logger.LevelInfo, "test")

		require.Len(t, buf, 1)
		assert.Equal(t, slog.LevelError, buf[0].Level)
		assert.Equal(t, "test", buf[0].Message)
	})

	t.Run("custom mapper preserved in derived adapters", func(t *testing.T) {
		t.Parallel()

		buf := entriesBuffer{}

		// Custom mapper that maps everything to WarnLevel
		customMapper := func(level int) slog.Level {
			return slog.LevelWarn
		}

		adapter := newAdapter(&buf, slogadapter.WithLogLevelMapper(customMapper))
		derived := adapter.WithFields(fields.F("foo", "bar"))

		// Both should use custom mapper
		adapter.Log(logger.LevelInfo, "original")
		derived.Log(logger.LevelDebug, "derived")

		require.Len(t, buf, 2)
		assert.Equal(t, slog.LevelWarn, buf[0].Level)
		assert.Equal(t, slog.LevelWarn, buf[1].Level)
		assert.Equal(t, "bar", buf[1].Attrs["foo"])
	})
}
