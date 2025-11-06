package zapadapter

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
)

type Option func(*Adapter)

// LogLevelMapper is a function that maps logger levels from
// [logger.Logger] to [zapcore.Level].
type LogLevelMapper func(level int) zapcore.Level

// WithLogLevelMapper sets custom log level mapper for the adapter.
func WithLogLevelMapper(fn LogLevelMapper) Option {
	return func(a *Adapter) {
		a.lvlMapper = fn
	}
}

// DefaultLogLevelMapper is a default implementation of [LogLevelMapper] that
// maps log levels from [logger.Logger] to appropriate [zapcore.Level].
//
// Note that zap does not have a separate trace level, so both LevelDebug and
// LevelTrace map to zapcore.DebugLevel.
func DefaultLogLevelMapper(level int) zapcore.Level {
	switch level {
	case logger.LevelError:
		return zapcore.ErrorLevel
	case logger.LevelWarning:
		return zapcore.WarnLevel
	case logger.LevelDebug,
		logger.LevelTrace:
		return zapcore.DebugLevel

	default:
		return zapcore.InfoLevel
	}
}

// Adapter of zap logger for [logger.Logger].
//
// This adapter guarantees support of stock logger's levels.
type Adapter struct {
	lgr *zap.Logger

	lvlMapper LogLevelMapper `exhaustruct:"optional"`
}

// New creates new logging adapter using provided [zap.Logger].
func New(lgr *zap.Logger, opts ...Option) *Adapter {
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
	if ce := a.lgr.Check(a.lvlMapper(level), msg); ce != nil {
		ce.Write(fieldsListToZapFields(fs)...)
	}
}

func (a *Adapter) WithFields(fs ...fields.Field) logger.Adapter {
	return &Adapter{
		lgr:       a.lgr.With(fieldsListToZapFields(fs)...),
		lvlMapper: a.lvlMapper,
	}
}

func (a *Adapter) Flush() error {
	return a.lgr.Sync() //nolint:wrapcheck
}

func fieldsListToZapFields(fs fields.List) []zap.Field {
	zfs := make([]zap.Field, 0, len(fs))

	for _, f := range fs {
		zfs = append(zfs, zap.Any(f.K, f.V))
	}

	return zfs
}
