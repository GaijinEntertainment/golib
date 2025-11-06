package logrusadapter

import (
	"github.com/sirupsen/logrus"

	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
)

type Option func(*Adapter)

// LogLevelMapper is a function that maps logger levels from
// [logger.Logger] to [logrus.Level].
type LogLevelMapper func(level int) logrus.Level

// WithLogLevelMapper sets custom log level mapper for the adapter.
func WithLogLevelMapper(fn LogLevelMapper) Option {
	return func(a *Adapter) {
		a.lvlMapper = fn
	}
}

// DefaultLogLevelMapper is a default implementation of [LogLevelMapper] that
// maps log levels from [logger.Logger] to appropriate [logrus.Level].
func DefaultLogLevelMapper(level int) logrus.Level {
	switch level {
	case logger.LevelError:
		return logrus.ErrorLevel
	case logger.LevelWarning:
		return logrus.WarnLevel
	case logger.LevelDebug:
		return logrus.DebugLevel
	case logger.LevelTrace:
		return logrus.TraceLevel

	default:
		return logrus.InfoLevel
	}
}

// Adapter of logrus logger for [logger.Logger].
//
// This adapter guarantees support of stock logger's levels.
type Adapter struct {
	lgr *logrus.Entry

	lvlMapper LogLevelMapper `exhaustruct:"optional"`
}

// New creates new logging adapter using provided [logrus.Entry].
func New(lgr *logrus.Entry, opts ...Option) *Adapter {
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

// Log implements [logger.Adapter.Log].
func (a *Adapter) Log(level int, msg string, fs ...fields.Field) {
	a.lgr.WithFields(fieldsListToLogrusFields(fs)).Log(a.lvlMapper(level), msg)
}

// WithFields implements [logger.Adapter.WithFields].
func (a *Adapter) WithFields(fs ...fields.Field) logger.Adapter {
	return &Adapter{
		lgr:       a.lgr.WithFields(fieldsListToLogrusFields(fs)),
		lvlMapper: a.lvlMapper,
	}
}

// Flush implements [logger.Adapter.Flush].
func (*Adapter) Flush() error {
	return nil
}

func fieldsListToLogrusFields(fs fields.List) logrus.Fields {
	lfs := make(logrus.Fields, len(fs))

	for _, f := range fs {
		lfs[f.K] = f.V
	}

	return lfs
}
