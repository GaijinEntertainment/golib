package slogadapter

import (
	"context"
	"log/slog"

	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
)

type Option func(*Adapter)

// LogLevelMapper is a function that maps logger levels from
// [logger.Logger] to [slog.Level].
type LogLevelMapper func(level int) slog.Level

// WithLogLevelMapper sets custom log level mapper for the adapter.
func WithLogLevelMapper(fn LogLevelMapper) Option {
	return func(a *Adapter) {
		a.lvlMapper = fn
	}
}

// DefaultLogLevelMapper is a default implementation of [LogLevelMapper] that
// maps log levels from [logger.Logger] to appropriate [slog.Level].
//
// Note that slog does not have a separate trace level, so both LevelDebug and
// LevelTrace map to slog.LevelDebug.
func DefaultLogLevelMapper(level int) slog.Level {
	switch level {
	case logger.LevelError:
		return slog.LevelError
	case logger.LevelWarning:
		return slog.LevelWarn
	case logger.LevelDebug,
		logger.LevelTrace:
		return slog.LevelDebug

	default:
		return slog.LevelInfo
	}
}

// Adapter of slog logger for [logger.Logger].
//
// This adapter guarantees support of stock logger's levels.
type Adapter struct {
	lgr *slog.Logger

	lvlMapper LogLevelMapper `exhaustruct:"optional"`
}

// New creates new logging adapter using provided [slog.Logger].
func New(lgr *slog.Logger, opts ...Option) *Adapter {
	a := &Adapter{
		lgr: lgr,
	}

	for _, opt := range opts {
		opt(a)
	}

	if a.lvlMapper == nil {
		a.lvlMapper = DefaultLogLevelMapper
	}

	return a
}

func (a *Adapter) Log(level int, msg string, fs ...fields.Field) {
	a.lgr.LogAttrs(context.Background(), a.lvlMapper(level), msg, fieldsListToSlogAttrs(fs)...)
}

func (a *Adapter) WithFields(fs ...fields.Field) logger.Adapter {
	return &Adapter{
		lgr:       slog.New(a.lgr.Handler().WithAttrs(fieldsListToSlogAttrs(fs))),
		lvlMapper: a.lvlMapper,
	}
}

func (*Adapter) Flush() error {
	return nil
}

func fieldsListToSlogAttrs(fs fields.List) []slog.Attr {
	zfs := make([]slog.Attr, 0, len(fs))

	for _, f := range fs {
		zfs = append(zfs, slog.Any(f.K, f.V))
	}

	return zfs
}
