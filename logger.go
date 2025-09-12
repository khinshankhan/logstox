package logstox

import (
	"github.com/khinshankhan/logstox/fields"
	"io"
)

// Logger is the small, portable logging interface.
// Structured data is passed as Field values constructed in field.go.
type Logger interface {
	// DEBUG (-1): for recording messages useful for debugging.
	Debug(string, ...fields.Field)
	// INFO (0): for messages describing normal application operations.
	Info(string, ...fields.Field)
	// WARN (1): for recording messages indicating something unusual happened that may need attention before it escalates to a more severe issue.
	Warn(string, ...fields.Field)
	// ERROR (2): for recording unexpected error conditions in the program.
	Error(string, ...fields.Field)
	// DPANIC (3): for recording severe error conditions in development. It behaves like PANIC in development and ERROR in production.
	DPanic(string, ...fields.Field)
	// PANIC (4): calls panic() after logging an error condition.
	Panic(string, ...fields.Field)
	// FATAL (5): calls os.Exit(1) after logging an error condition.
	Fatal(string, ...fields.Field)

	// With creates a child logger and adds structured context to it. Fields added
	// to the child don't affect the parent, and vice versa. Any fields that
	// require evaluation (such as Objects) are evaluated upon invocation of With.
	With(...fields.Field) Logger
	// Named adds a new path segment to the logger's name. Segments are joined by
	// periods. By default, Loggers are unnamed.
	Named(string) Logger
	// Sync calls the underlying Core's Sync method, flushing any buffered log
	// entries. Applications should take care to call Sync before exiting.
	Sync() error
}

// LevelCheck is an optional extension that reports if a level is enabled.
// Backends that can answer cheaply may implement this.
type LevelCheck interface {
	Enabled(Level) bool
}

// Options are hints used by a Backend when constructing a Logger.
// Backends may choose to ignore some fields.
type Options struct {
	Level      Level          // minimum level to record
	AddSource  bool           // include file:line when supported
	Name       string         // initial logger scope
	Writer     io.Writer      // preferred sink (backend may ignore)
	TimeLayout string         // eg time.RFC3339Nano (backend may ignore)
	Fields     []fields.Field // default fields for the base logger
}

// Backend builds a Logger from Options.
// Backends live in subpackages (eg backend/zapx, backend/slogx).
// Or consumers roll out their custom backend.
type Backend interface {
	New(Options) Logger
}
