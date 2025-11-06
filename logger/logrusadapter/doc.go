// Package logrusadapter provides a logger adapter for the logrus logging library.
//
// This adapter enables using logrus as a backend for logger.Logger, bridging
// the logrus API with the logger abstraction. For documentation on the logger
// API itself, see the [dev.gaijin.team/go/golib/logger] package.
//
// # Basic Usage
//
// Create a logrus logger, wrap it in an adapter, and use with logger.Logger:
//
//	import (
//		"github.com/sirupsen/logrus"
//		"dev.gaijin.team/go/golib/logger"
//		"dev.gaijin.team/go/golib/logger/logrusadapter"
//	)
//
//	ll := logrus.New()
//	ll.SetFormatter(&logrus.JSONFormatter{})
//	ll.SetLevel(logrus.TraceLevel)
//
//	adapter := logrusadapter.New(logrus.NewEntry(ll))
//	log := logger.New(adapter, logger.LevelInfo)
//
//	log.Info("server started", fields.F("port", 8080))
//
// IMPORTANT: Configure logrus with the highest level (TraceLevel) to ensure all
// logs are passed to the adapter. Level filtering is performed by logger.Logger,
// not by logrus. If logrus is set to a lower level (e.g., InfoLevel), it will
// filter out debug and trace logs if they would be sent by logger.Logger.
//
// # Custom Level Mapping
//
// By default, logger levels map to logrus levels using [DefaultLogLevelMapper].
// You can provide a custom mapper if you need different level mapping behavior:
//
//	customMapper := func(level int) logrus.Level {
//		switch level {
//		case logger.LevelError:
//			return logrus.FatalLevel  // Map errors to fatal
//		default:
//			return logrus.InfoLevel
//		}
//	}
//
//	adapter := logrusadapter.New(
//		logrus.NewEntry(ll),
//		logrusadapter.WithLogLevelMapper(customMapper),
//	)
package logrusadapter
