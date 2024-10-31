package zapadapter

import (
	"go.uber.org/zap"

	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
)

type zapF func(*zap.Logger, string, ...zap.Field)
type logMapper [logger.LevelTrace + 1]zapF

var mapper logMapper //nolint:gochecknoglobals

func init() { //nolint:gochecknoinits
	mapper[logger.LevelError] = (*zap.Logger).Error
	mapper[logger.LevelWarning] = (*zap.Logger).Warn
	mapper[logger.LevelInfo] = (*zap.Logger).Info
	mapper[logger.LevelDebug] = (*zap.Logger).Debug
	mapper[logger.LevelTrace] = (*zap.Logger).Debug
}

// Adapter of zap logger for [logger.Logger].
//
// This adapter guarantees support of stock logger's levels.
type Adapter struct {
	lgr *zap.Logger
}

// New creates new logging adapter using provided [zap.Logger].
//
// Note, that by the contract of [logger.Logger], adapter should not perform
// level-filtering internally - it is done by [logger.Logger] itself.
func New(lgr *zap.Logger) *Adapter {
	return &Adapter{
		lgr: lgr,
	}
}

func (a *Adapter) Log(level int, msg string, err error, fs ...fields.Field) {
	var zf zapF
	if level < len(mapper) {
		zf = mapper[level]
	}

	if zf == nil {
		zf = (*zap.Logger).Info

		a.lgr.Error("Unknown log level", zap.Int("got-level", level))
	}

	zf(a.lgr, msg, fieldsListToZapFields(fs, err)...)
}

func (a *Adapter) WithFields(fs ...fields.Field) logger.Adapter {
	return &Adapter{
		lgr: a.lgr.With(fieldsListToZapFields(fs, nil)...),
	}
}

func (a *Adapter) WithName(name string) logger.Adapter {
	return &Adapter{
		lgr: a.lgr.Named(name),
	}
}

func (a *Adapter) WithStackTrace(_ string) logger.Adapter {
	return a
}

func (a *Adapter) Flush() error {
	return a.lgr.Sync() //nolint:wrapcheck
}

func fieldsListToZapFields(fs fields.List, err error) []zap.Field {
	zfs := make([]zap.Field, 0, len(fs)+boolToInt(err != nil))

	if err != nil {
		zfs = append(zfs, zap.Error(err))
	}

	for _, f := range fs {
		zfs = append(zfs, zap.Any(f.K, f.V))
	}

	return zfs
}

func boolToInt(v bool) int {
	if v {
		return 1
	}

	return 0
}
