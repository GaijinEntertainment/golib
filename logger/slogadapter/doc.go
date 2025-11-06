// Package slogadapter provides a logger adapter for the standard library's slog
// logging package.
//
// This adapter enables using slog as a backend for [logger.Logger], bridging the
// slog API with the logger abstraction. For documentation on the logger API
// itself, see the [dev.gaijin.team/go/golib/logger] package.
//
// # Basic Usage
//
// Create a slog logger, wrap it in an adapter, and use with logger.Logger:
//
//	import (
//		"log/slog"
//		"os"
//		"dev.gaijin.team/go/golib/logger"
//		"dev.gaijin.team/go/golib/logger/slogadapter"
//	)
//
//	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
//		Level: slog.LevelDebug,
//	})
//	sl := slog.New(handler)
//
//	adapter := slogadapter.New(sl)
//	log := logger.New(adapter, logger.LevelInfo)
//
//	log.Info("server started", fields.F("port", 8080))
//
// IMPORTANT: Configure slog handler with the lowest level (LevelDebug or lower)
// to ensure all logs are passed to the adapter. Level filtering is performed by
// logger.Logger, not by slog. If slog handler is set to a higher level (e.g.,
// LevelWarn), it will filter out info and debug logs if they would be sent by
// logger.Logger.
//
// # Custom Level Mapping
//
// By default, logger levels map to slog levels using [DefaultLogLevelMapper].
// You can provide a custom mapper if you need different level mapping behavior:
//
//	customMapper := func(level int) slog.Level {
//		switch level {
//		case logger.LevelError:
//			return slog.LevelError + 4  // Map errors to higher severity
//		default:
//			return slog.LevelInfo
//		}
//	}
//
//	adapter := slogadapter.New(
//		sl,
//		slogadapter.WithLogLevelMapper(customMapper),
//	)
//
// Note that slog does not have a separate trace level, so both LevelDebug and
// LevelTrace map to slog.LevelDebug by default.
package slogadapter
