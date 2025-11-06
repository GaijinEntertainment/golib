// Package zapadapter provides a logger adapter for the zap logging library.
//
// This adapter enables using zap as a backend for [logger.Logger], bridging the
// zap API with the logger abstraction. For documentation on the logger API
// itself, see the [dev.gaijin.team/go/golib/logger] package.
//
// # Basic Usage
//
// Create a zap logger, wrap it in an adapter, and use with logger.Logger:
//
//	import (
//		"go.uber.org/zap"
//		"go.uber.org/zap/zapcore"
//		"dev.gaijin.team/go/golib/logger"
//		"dev.gaijin.team/go/golib/logger/zapadapter"
//	)
//
//	config := zap.NewProductionConfig()
//	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
//	zl, _ := config.Build()
//
//	adapter := zapadapter.New(zl)
//	log := logger.New(adapter, logger.LevelInfo)
//
//	log.Info("server started", fields.F("port", 8080))
//
// IMPORTANT: Configure zap with the lowest level (DebugLevel or lower) to ensure
// all logs are passed to the adapter. Level filtering is performed by
// logger.Logger, not by zap. If zap is set to a higher level (e.g., WarnLevel),
// it will filter out info and debug logs if they would be sent by logger.Logger.
//
// # Custom Level Mapping
//
// By default, logger levels map to zap levels using [DefaultLogLevelMapper]. You
// can provide a custom mapper if you need different level mapping behavior:
//
//	customMapper := func(level int) zapcore.Level {
//		switch level {
//		case logger.LevelError:
//			return zapcore.ErrorLevel + 1  // Map errors to higher severity
//		default:
//			return zapcore.InfoLevel
//		}
//	}
//
//	adapter := zapadapter.New(
//		zl,
//		zapadapter.WithLogLevelMapper(customMapper),
//	)
//
// Note that zap does not have a separate trace level, so both LevelDebug and
// LevelTrace map to zapcore.DebugLevel by default.
package zapadapter
