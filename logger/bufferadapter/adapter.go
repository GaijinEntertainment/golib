package bufferadapter

import (
	"slices"

	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
)

type LogEntry struct {
	LoggerName string `exhaustruct:"optional"`
	Level      int
	Msg        string
	Error      error       `exhaustruct:"optional"`
	Fields     fields.List `exhaustruct:"optional"`
}

type LogEntries []LogEntry

func (le *LogEntries) Reset() {
	*le = LogEntries{}
}

// Adapter is a logger.Adapter that writes logs to the provided [LogEntries]
// slice.
//
// This adapter is designed for testing purposes and is not intended for
// production use.
type Adapter struct {
	buff *LogEntries
	name string      `exhaustruct:"optional"`
	fs   fields.List `exhaustruct:"optional"`
}

// New creates a new [Adapter] instance.
func New(buff *LogEntries) *Adapter {
	return &Adapter{buff: buff}
}

func (a *Adapter) Log(level int, msg string, err error, fs ...fields.Field) {
	e := LogEntry{
		LoggerName: a.name,
		Level:      level,
		Msg:        msg,
		Error:      err,
	}

	e.Fields = append(slices.Clone(a.fs), fs...)

	*a.buff = append(*a.buff, e)
}

func (a *Adapter) WithFields(fs ...fields.Field) logger.Adapter {
	return &Adapter{
		buff: a.buff,
		name: a.name,
		fs:   append(slices.Clone(a.fs), fs...),
	}
}

func (a *Adapter) WithName(name string) logger.Adapter {
	return &Adapter{
		buff: a.buff,
		name: name,
		fs:   slices.Clone(a.fs),
	}
}

func (a *Adapter) WithStackTrace(_ string) logger.Adapter {
	return &Adapter{
		buff: a.buff,
		name: a.name,
		fs:   slices.Clone(a.fs),
	}
}

func (*Adapter) Flush() error {
	return nil
}
