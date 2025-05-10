package e_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"dev.gaijin.team/go/golib/e"
	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
	"dev.gaijin.team/go/golib/logger/bufferadapter"
)

func TestLog(t *testing.T) {
	t.Parallel()

	entries := make(bufferadapter.LogEntries, 0)
	log := logger.New(bufferadapter.New(&entries), logger.DefaultLogLevel)

	e1 := errors.New("e1") //nolint:err113
	e2 := fmt.Errorf("e2: %w", e1)
	e3 := e.NewFrom("e3", e1, fields.F("foo", "bar"))
	e4 := e.New("e4")

	e.Log(nil, log.Error)
	e.Log(e1, log.Error)
	e.Log(e1, log.WarningE)
	e.Log(e2, log.InfoE)
	e.Log(e3, log.DebugE)
	e.Log(e4, log.TraceE)

	assert.Equal(t, bufferadapter.LogEntries{
		{
			LoggerName: "",
			Level:      logger.LevelError,
			Msg:        e1.Error(),
			Error:      nil,
			Fields:     nil,
		},
		{
			LoggerName: "",
			Level:      logger.LevelWarning,
			Msg:        e1.Error(),
			Error:      nil,
			Fields:     nil,
		},
		{
			LoggerName: "",
			Level:      logger.LevelInfo,
			Msg:        e2.Error(),
			Error:      nil,
			Fields:     nil,
		},
		{
			LoggerName: "",
			Level:      logger.LevelDebug,
			Msg:        e3.Reason(),
			Error:      e1,
			Fields:     e3.Fields(),
		},
		{
			LoggerName: "",
			Level:      logger.LevelTrace,
			Msg:        e4.Reason(),
			Error:      nil,
			Fields:     nil,
		},
	}, entries)
}
