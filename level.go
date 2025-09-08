package logstox

import (
	"encoding"
	"fmt"
	"strings"
)

// TODO: migrate to using a tool like stringer to generate interface implementations.

// Level is an implementation-agnostic log severity.
// Higher numbers are more severe. The zero value is InfoLevel.
type Level int8

const (
	// DebugLevel enables verbose diagnostic output, usually disabled in production.
	DebugLevel Level = -1
	// InfoLevel describes normal application operations (default).
	InfoLevel Level = 0
	// WarnLevel indicates unusual conditions that may need attention, but don't need individual human review.
	WarnLevel Level = 1
	// ErrorLevel records unexpected errors. If an application is running smoothly, it shouldn't generate any
	// error-level logs.
	ErrorLevel Level = 2
	// DPanicLevel is intended for development-only severe errors
	// (often treated like Panic in dev and Error in prod by loggers).
	DPanicLevel Level = 3
	// PanicLevel is intended for errors after which a logger will panic.
	PanicLevel Level = 4
	// FatalLevel is intended for errors after which a logger will exit.
	FatalLevel Level = 5
)

// String implements fmt.Stringer, returning the lower-case name of the log level.
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case DPanicLevel:
		return "dpanic"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", int(l))
	}
}

// Valid reports whether l is one of the defined levels.
func (l Level) Valid() bool {
	switch l {
	case DebugLevel, InfoLevel, WarnLevel, ErrorLevel, DPanicLevel, PanicLevel, FatalLevel:
		return true
	default:
		return false
	}
}

// ParseLevel parses s (case-insensitive, leading/trailing spaces ignored) into a Level.
func ParseLevel(s string) (Level, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return DebugLevel, nil
	case "info":
		return InfoLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "error":
		return ErrorLevel, nil
	case "dpanic":
		return DPanicLevel, nil
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	default:
		return 0, fmt.Errorf("unknown level %q", s)
	}
}

// MarshalText implements encoding.TextMarshaler.
func (l Level) MarshalText() ([]byte, error) {
	if !l.Valid() {
		return nil, fmt.Errorf("cannot marshal invalid level %d", int(l))
	}
	return []byte(l.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (l *Level) UnmarshalText(text []byte) error {
	parsed, err := ParseLevel(string(text))
	if err != nil {
		return err
	}
	*l = parsed
	return nil
}

// Interface satisfaction (compile-time assertions).
var (
	_ fmt.Stringer             = (*Level)(nil)
	_ encoding.TextMarshaler   = (*Level)(nil)
	_ encoding.TextUnmarshaler = (*Level)(nil)
)
