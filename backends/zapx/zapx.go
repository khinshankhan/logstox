package zapx

import (
	"time"

	"github.com/khinshankhan/logstox"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Backend implements logstox.Backend using Uber's zap.
type Backend struct {
	// If true, start from zap.NewDevelopmentConfig(); else production config.
	Development bool
	// Optional layout for timestamps (defaults to time.RFC3339Nano).
	TimeLayout string
	// If true, include file:line via zap.AddCaller().
	AddSource bool
}

// Interface satisfaction (compile-time assertions).
var (
	_ logstox.Backend = Backend{}
)

// New builds a zap-backed logger from Options.
func (b Backend) New(opts logstox.Options) logstox.Logger {
	var cfg zap.Config
	if b.Development {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	encCfg := cfg.EncoderConfig
	layout, _ := firstNonEmptyString(
		opts.TimeLayout,
		b.TimeLayout,
		time.RFC3339Nano,
	)
	encCfg.EncodeTime = zapcore.TimeEncoderOfLayout(layout)

	// Keep stacktraces out unless explicitly added via Native(zap.Stack(...)).
	encCfg.StacktraceKey = ""
	cfg.EncoderConfig = encCfg

	// Level override if provided?
	if zl, ok := toZapLevel(opts.Level); ok {
		cfg.Level = zap.NewAtomicLevelAt(zl)
	}

	// core logger
	var base *zap.Logger
	if opts.Writer != nil {
		enc := zapcore.NewJSONEncoder(encCfg)
		ws := zapcore.AddSync(opts.Writer)
		core := zapcore.NewCore(enc, ws, cfg.Level)
		var optsZap []zap.Option
		if b.AddSource || opts.AddSource {
			optsZap = append(optsZap, zap.AddCaller(), zap.AddCallerSkip(1))
		}
		base = zap.New(core, optsZap...)
	} else {
		var optsZap []zap.Option
		if b.AddSource || opts.AddSource {
			optsZap = append(optsZap, zap.AddCaller(), zap.AddCallerSkip(1))
		}
		base = zap.Must(cfg.Build(optsZap...))
	}

	if opts.Name != "" {
		base = base.Named(opts.Name)
	}
	if len(opts.Fields) > 0 {
		base = base.With(toZapFields(base, zapcore.InfoLevel, opts.Fields...)...)
	}

	return &zlogger{l: base}
}

// zlogger is a zap-backed implementation of logstox.Logger
type zlogger struct{ l *zap.Logger }

// Interface satisfaction (compile-time assertions).
var (
	_ logstox.Logger = (*zlogger)(nil)
)

func (lg *zlogger) Debug(m string, f ...logstox.Field) {
	lg.l.Debug(m, toZapFields(lg.l, zapcore.DebugLevel, f...)...)
}
func (lg *zlogger) Info(m string, f ...logstox.Field) {
	lg.l.Info(m, toZapFields(lg.l, zapcore.InfoLevel, f...)...)
}
func (lg *zlogger) Warn(m string, f ...logstox.Field) {
	lg.l.Warn(m, toZapFields(lg.l, zapcore.WarnLevel, f...)...)
}
func (lg *zlogger) Error(m string, f ...logstox.Field) {
	lg.l.Error(m, toZapFields(lg.l, zapcore.ErrorLevel, f...)...)
}
func (lg *zlogger) DPanic(m string, f ...logstox.Field) {
	lg.l.DPanic(m, toZapFields(lg.l, zapcore.DPanicLevel, f...)...)
}
func (lg *zlogger) Panic(m string, f ...logstox.Field) {
	lg.l.Panic(m, toZapFields(lg.l, zapcore.PanicLevel, f...)...)
}
func (lg *zlogger) Fatal(m string, f ...logstox.Field) {
	lg.l.Fatal(m, toZapFields(lg.l, zapcore.FatalLevel, f...)...)
}

// Optional level check (handy for guarding expensive field prep).
func (lg *zlogger) Enabled(l logstox.Level) bool {
	zl, ok := toZapLevel(l)
	if !ok {
		return lg.l.Core().Enabled(zapcore.InfoLevel)
	}
	return lg.l.Core().Enabled(zl)
}

func (lg *zlogger) With(f ...logstox.Field) logstox.Logger {
	return &zlogger{l: lg.l.With(toZapFields(lg.l, zapcore.InfoLevel, f...)...)}
}
func (lg *zlogger) Named(n string) logstox.Logger { return &zlogger{l: lg.l.Named(n)} }
func (lg *zlogger) Sync() error                   { return lg.l.Sync() }
