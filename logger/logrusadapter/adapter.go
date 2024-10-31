package logrusadapter

import (
	"github.com/sirupsen/logrus"

	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
)

type logrusF func(*logrus.Entry, ...any)
type logMapper [logger.LevelTrace + 1]logrusF

//nolint:gochecknoglobals
var mapper logMapper

//nolint:gochecknoinits
func init() {
	mapper[logger.LevelError] = (*logrus.Entry).Error
	mapper[logger.LevelWarning] = (*logrus.Entry).Warn
	mapper[logger.LevelInfo] = (*logrus.Entry).Info
	mapper[logger.LevelDebug] = (*logrus.Entry).Debug
	mapper[logger.LevelTrace] = (*logrus.Entry).Trace
}

const LoggerNameKey = "logger"

// Adapter of logrus logger for [logger.Logger].
//
// This adapter guarantees support of stock logger's levels.
type Adapter struct {
	lgr *logrus.Entry
}

// New creates new logging adapter using provided [logrus.Entry].
//
// Note, that by the contract of [logger.Logger], adapter should not perform
// level-filtering internally - it is done by [logger.Logger] itself.
func New(lgr *logrus.Entry) *Adapter {
	return &Adapter{
		lgr: lgr,
	}
}

func (a *Adapter) Log(level int, msg string, err error, fs ...fields.Field) {
	lgr := a.lgr

	var lf logrusF

	if level < len(mapper) {
		lf = mapper[level]
	}

	if lf == nil {
		lf = (*logrus.Entry).Info

		lgr.WithField("got-level", level).Error("Unknown log level")
	}

	if err != nil {
		lgr = lgr.WithError(err)
	}

	lf(lgr.WithFields(fieldsListToLogrusFields(fs)), msg)
}

func (a *Adapter) WithFields(fs ...fields.Field) logger.Adapter {
	return &Adapter{
		lgr: a.lgr.WithFields(fieldsListToLogrusFields(fs)),
	}
}

// WithName returns a logger adapter with the given name attached to it. As
// logrus does not support logger names, this adds [LoggerNameKey] field.
func (a *Adapter) WithName(name string) logger.Adapter {
	return &Adapter{
		lgr: a.lgr.WithField(LoggerNameKey, name),
	}
}

func (a *Adapter) WithStackTrace(_ string) logger.Adapter {
	return a
}

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
