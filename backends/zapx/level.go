package zapx

import (
	"github.com/khinshankhan/logstox"

	"go.uber.org/zap/zapcore"
)

func toZapLevel(l logstox.Level) (zapcore.Level, bool) {
	switch l {
	case logstox.DebugLevel:
		return zapcore.DebugLevel, true
	case logstox.InfoLevel:
		return zapcore.InfoLevel, true
	case logstox.WarnLevel:
		return zapcore.WarnLevel, true
	case logstox.ErrorLevel:
		return zapcore.ErrorLevel, true
	case logstox.DPanicLevel:
		return zapcore.DPanicLevel, true
	case logstox.PanicLevel:
		return zapcore.PanicLevel, true
	case logstox.FatalLevel:
		return zapcore.FatalLevel, true
	default:
		return zapcore.InfoLevel, false
	}
}
