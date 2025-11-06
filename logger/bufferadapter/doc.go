// Package bufferadapter provides an in-memory logger adapter for testing.
//
// This package implements the logger.Adapter interface by storing log entries in
// memory, allowing tests to inspect what was logged. The adapter is safe for
// concurrent use and designed specifically for testing purposes.
//
// # Usage
//
// Create a buffer adapter and use it with a logger:
//
//	adapter, buff := bufferadapter.New()
//	log := logger.New(adapter, logger.LevelInfo)
//
//	log.Info("test message", fields.F("key", "value"))
//
//	// Assert in tests
//	entries := buff.GetAll()
//	require.Len(t, entries, 1)
//	assert.Equal(t, "test message", entries[0].Msg)
//	assert.Equal(t, logger.LevelInfo, entries[0].Level)
package bufferadapter
