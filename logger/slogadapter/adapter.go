package slogadapter

import (
	"context"
	"log/slog"

	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
)

// Adapter of slog logger for [logger.Logger].
//
// This adapter guarantees support of stock logger's levels.
type Adapter struct {
	lgr *slog.Logger
}

// New creates new logging adapter using provided [slog.Logger].
//
// Note, that by the contract of [logger.Logger], adapter should not perform
// level-filtering internally - it is done by [logger.Logger] itself.
func New(lgr *slog.Logger) *Adapter {
	return &Adapter{
		lgr: lgr,
	}
}

func (a *Adapter) Log(level int, msg string, err error, fs ...fields.Field) {
	sl := slog.LevelInfo

	switch level {
	case logger.LevelError:
		sl = slog.LevelError

	case logger.LevelWarning:
		sl = slog.LevelWarn

	case logger.LevelInfo:
		sl = slog.LevelInfo

	case logger.LevelDebug:
		sl = slog.LevelDebug

	case logger.LevelTrace:
		sl = slog.LevelDebug

	default:
		a.lgr.Error("Unknown log level", slog.Int("got-level", level))
	}

	a.lgr.LogAttrs(context.Background(), sl, msg, fieldsListToSlogAttrs(fs, err)...)
}

func (a *Adapter) WithFields(fs ...fields.Field) logger.Adapter {
	return &Adapter{
		lgr: slog.New(a.lgr.Handler().WithAttrs(fieldsListToSlogAttrs(fs, nil))),
	}
}

func (a *Adapter) WithName(name string) logger.Adapter {
	return &Adapter{
		lgr: a.lgr.WithGroup(name),
	}
}

func (a *Adapter) WithStackTrace(_ string) logger.Adapter {
	return a
}

func (*Adapter) Flush() error {
	return nil
}

func fieldsListToSlogAttrs(fs fields.List, err error) []slog.Attr {
	zfs := make([]slog.Attr, 0, len(fs)+boolToInt(err != nil))

	if err != nil {
		zfs = append(zfs, slog.Any("error", err))
	}

	for _, f := range fs {
		zfs = append(zfs, slog.Any(f.K, f.V))
	}

	return zfs
}

func boolToInt(v bool) int {
	if v {
		return 1
	}

	return 0
}
