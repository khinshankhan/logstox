package zapx

import (
	"time"

	"github.com/khinshankhan/logstox"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapField = zap.Field

// Backend builds a Logger[ZapField] from logstox.Options[ZapField].
type Backend struct {
	Development bool
	// If non-empty, sets the timestamp layout (eg, time.RFC3339Nano).
	// Options.TimeLayout takes precedence over this.
	TimeLayout string
	AddSource  bool
	// CallerSkip controls the number of stack frames to skip when reporting the caller.
	// If zero, defaults to 1, which typically points to the caller of the logger method.
	// Increase this value if wrapping the logger in additional abstraction layers.
	CallerSkip int
}

// Interface satisfaction (compile-time assertions).
var _ logstox.Backend[ZapField] = Backend{}

// firstNonEmpty is a quick utility function to choose between provided options or fall back to a default
func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

// New constructs a zap-backed Logger[ZapField].
func (b Backend) New(o logstox.Options[ZapField]) logstox.Logger[ZapField] {
	// Base config: dev/prod
	var cfg zap.Config
	if b.Development {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	// Encoder config, time layout + hide stacktrace unless explicitly added
	enc := cfg.EncoderConfig
	layout := firstNonEmpty(o.TimeLayout, b.TimeLayout, time.RFC3339Nano)
	enc.EncodeTime = zapcore.TimeEncoderOfLayout(layout)
	enc.StacktraceKey = ""
	cfg.EncoderConfig = enc

	// Level override from Options if provided/ mapped.
	if zl, ok := toZapLevel(o.Level); ok {
		cfg.Level = zap.NewAtomicLevelAt(zl)
	}

	// Build logger, either via provided Writer or default sinks.
	var base *zap.Logger
	var opts []zap.Option
	if b.AddSource || o.AddSource {
		skip := b.CallerSkip
		if skip == 0 {
			skip = 1
		}
		// AddCallerSkip(2) to point at the user's callsite (skipping our wrapper method).
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(skip))
	}

	if o.Writer != nil {
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(enc),
			zapcore.AddSync(o.Writer),
			cfg.Level,
		)
		base = zap.New(core, opts...)
	} else {
		base = zap.Must(cfg.Build(opts...))
	}

	if o.Name != "" {
		base = base.Named(o.Name)
	}
	if len(o.Fields) > 0 {
		base = base.With(o.Fields...)
	}

	return logger{l: base}
}

// logger is a thin zap-backed implementation of logstox.Logger[ZapField].
type logger struct{ l *zap.Logger }

// Interface satisfaction (compile-time assertions).
var _ logstox.Logger[ZapField] = logger{}

// DEBUG (-1): for recording messages useful for debugging.
func (lg logger) Debug(m string, f ...ZapField) { lg.l.Debug(m, f...) }

// INFO (0): for messages describing normal application operations.
func (lg logger) Info(m string, f ...ZapField) { lg.l.Info(m, f...) }

// WARN (1): for recording messages indicating something unusual happened that may need attention before it escalates to a more severe issue.
func (lg logger) Warn(m string, f ...ZapField) { lg.l.Warn(m, f...) }

// ERROR (2): for recording unexpected error conditions in the program.
func (lg logger) Error(m string, f ...ZapField) { lg.l.Error(m, f...) }

// DPANIC (3): for recording severe error conditions in development. It behaves like PANIC in development and ERROR in production.
func (lg logger) DPanic(m string, f ...ZapField) { lg.l.DPanic(m, f...) }

// PANIC (4): calls panic() after logging an error condition.
func (lg logger) Panic(m string, f ...ZapField) { lg.l.Panic(m, f...) }

// FATAL (5): calls os.Exit(1) after logging an error condition.
func (lg logger) Fatal(m string, f ...ZapField) { lg.l.Fatal(m, f...) }

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa. Any fields that
// require evaluation (such as Objects) are evaluated upon invocation of With.
func (lg logger) With(f ...ZapField) logstox.Logger[ZapField] {
	return logger{l: lg.l.With(f...)}
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func (lg logger) Named(n string) logstox.Logger[ZapField] {
	return logger{l: lg.l.Named(n)}
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func (lg logger) Sync() error {
	return lg.l.Sync()
}
