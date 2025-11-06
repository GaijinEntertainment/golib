package bufferadapter

import (
	"slices"
	"sync"

	"dev.gaijin.team/go/golib/fields"
	"dev.gaijin.team/go/golib/logger"
)

type LogEntry struct {
	Level  int
	Msg    string
	Fields fields.List `exhaustruct:"optional"`
}

// LogEntries is a buffer for storing log entries in memory. It is safe for
// concurrent use.
type LogEntries struct {
	mu      sync.RWMutex `exhaustruct:"optional"`
	entries []LogEntry   `exhaustruct:"optional"`
}

// Add appends a log entry.
func (le *LogEntries) Add(entry LogEntry) {
	le.mu.Lock()
	defer le.mu.Unlock()

	le.entries = append(le.entries, entry)
}

// Reset clears all log entries, preserving capacity.
func (le *LogEntries) Reset() {
	le.mu.Lock()
	defer le.mu.Unlock()

	le.entries = le.entries[:0]
}

// Len returns the number of log entries.
func (le *LogEntries) Len() int {
	le.mu.RLock()
	defer le.mu.RUnlock()
	return len(le.entries)
}

// Get returns a copy of the log entry at the given index.
// Panics if index is out of range.
func (le *LogEntries) Get(i int) LogEntry {
	le.mu.RLock()
	defer le.mu.RUnlock()
	return le.entries[i]
}

// GetAll returns a copy of log entries buffer.
func (le *LogEntries) GetAll() []LogEntry {
	le.mu.RLock()
	defer le.mu.RUnlock()
	return slices.Clone(le.entries)
}

// Adapter is a logger.Adapter implementation that writes logs to the provided
// [LogEntries] collection.
type Adapter struct {
	buff *LogEntries `exhaustruct:"optional"`
	fs   fields.List `exhaustruct:"optional"`
}

// New creates a new [Adapter] instance along with a new [LogEntries] buffer that
// will be used by adapter.
func New() (*Adapter, *LogEntries) {
	buff := &LogEntries{
		entries: make([]LogEntry, 0),
	}

	return &Adapter{buff: buff}, buff
}

func (a *Adapter) Log(level int, msg string, fs ...fields.Field) {
	e := LogEntry{
		Level:  level,
		Msg:    msg,
		Fields: append(slices.Clone(a.fs), fs...),
	}

	a.buff.Add(e)
}

func (a *Adapter) WithFields(fs ...fields.Field) logger.Adapter {
	return &Adapter{
		buff: a.buff,
		fs:   append(slices.Clone(a.fs), fs...),
	}
}

func (*Adapter) Flush() error {
	return nil
}
