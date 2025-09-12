package logstox

import (
	"io"
)

// Logger is the small, portable logging interface parameterized by the field type FT.
// Structured data is passed as Field values constructed in field.go.

type Logger[FT any] interface {
	// DEBUG (-1): for recording messages useful for debugging.
	Debug(string, ...FT)
	// INFO (0): for messages describing normal application operations.
	Info(string, ...FT)
	// WARN (1): for recording messages indicating something unusual happened that may need attention before it escalates to a more severe issue.
	Warn(string, ...FT)
	// ERROR (2): for recording unexpected error conditions in the program.
	Error(string, ...FT)
	// DPANIC (3): for recording severe error conditions in development. It behaves like PANIC in development and ERROR in production.
	DPanic(string, ...FT)
	// PANIC (4): calls panic() after logging an error condition.
	Panic(string, ...FT)
	// FATAL (5): calls os.Exit(1) after logging an error condition.
	Fatal(string, ...FT)

	// With creates a child logger and adds structured context to it. Fields added
	// to the child don't affect the parent, and vice versa. Any fields that
	// require evaluation (such as Objects) are evaluated upon invocation of With.
	With(...FT) Logger[FT]
	// Named adds a new path segment to the logger's name. Segments are joined by
	// periods. By default, Loggers are unnamed.
	Named(string) Logger[FT]
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
type Options[FT any] struct {
	Level      Level     // minimum level to record
	AddSource  bool      // include file:line when supported
	Name       string    // initial logger scope
	Writer     io.Writer // preferred sink (backend may ignore)
	TimeLayout string    // eg time.RFC3339Nano (backend may ignore)
	Fields     []FT      // default fields for the base logger
}

// Backend builds a Logger from Options all parameterized by the field type FT.
// Backends live in subpackages (eg backend/zapx, backend/slogx).
// Or consumers roll out their custom backend.
type Backend[FT any] interface {
	New(Options[FT]) Logger[FT]
}
