package zapx

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/khinshankhan/logstox"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func toZapFields(l *zap.Logger, lvl zapcore.Level, fs ...logstox.Field) []zap.Field {
	if len(fs) == 0 {
		return nil
	}
	out := make([]zap.Field, 0, len(fs))
	enabled := l.Core().Enabled(lvl)

	for _, f := range fs {
		if f.IsSkip() {
			continue
		}
		switch f.Kind() {

		// lazy
		case logstox.FieldKindLazyValue:
			if !enabled {
				continue
			}
			fn := f.Value.(func() []logstox.Field)
			sub := toZapFields(l, lvl, fn()...)
			out = append(out, sub...)

		// special
		case logstox.FieldKindDict:
			out = append(out, zap.Object(f.Key, dictMarshaler{fs: f.Value.([]logstox.Field)}))
		case logstox.FieldKindRawJSON:
			out = append(out, zap.Any(f.Key, json.RawMessage(f.Value.([]byte))))
		case logstox.FieldKindHexBytes:
			out = append(out, zap.String(f.Key, hex.EncodeToString(f.Value.([]byte))))
		case logstox.FieldKindTimestamp:
			t := f.Value.(time.Time)
			if t.IsZero() {
				t = time.Now()
			}
			out = append(out, zap.Time(f.Key, t))

		// scalars
		case logstox.FieldKindString:
			out = append(out, zap.String(f.Key, f.Value.(string)))
		case logstox.FieldKindBool:
			out = append(out, zap.Bool(f.Key, f.Value.(bool)))
		case logstox.FieldKindInt64:
			out = append(out, zap.Int64(f.Key, f.Value.(int64)))
		case logstox.FieldKindUint64:
			out = append(out, zap.Uint64(f.Key, f.Value.(uint64)))
		case logstox.FieldKindFloat64:
			out = append(out, zap.Float64(f.Key, f.Value.(float64)))
		case logstox.FieldKindDuration:
			out = append(out, zap.Duration(f.Key, f.Value.(time.Duration)))
		case logstox.FieldKindTime:
			out = append(out, zap.Time(f.Key, f.Value.(time.Time)))
		case logstox.FieldKindError:
			err := f.Value.(error)
			if f.Key == "" || f.Key == logstox.ErrorKey {
				out = append(out, zap.Error(err))
			} else {
				out = append(out, zap.NamedError(f.Key, err))
			}

		// slices
		case logstox.FieldKindStrings:
			out = append(out, zap.Strings(f.Key, f.Value.([]string)))
		case logstox.FieldKindBools:
			out = append(out, zap.Bools(f.Key, f.Value.([]bool)))
		case logstox.FieldKindInt64s:
			out = append(out, zap.Int64s(f.Key, f.Value.([]int64)))
		case logstox.FieldKindUint64s:
			out = append(out, zap.Uint64s(f.Key, f.Value.([]uint64)))
		case logstox.FieldKindFloat64s:
			out = append(out, zap.Float64s(f.Key, f.Value.([]float64)))
		case logstox.FieldKindErrors:
			out = append(out, zap.Errors(f.Key, f.Value.([]error)))

		// any / default
		case logstox.FieldKindAny:
			out = append(out, zap.Any(f.Key, f.Value))

		default:
			out = append(out, zap.Skip())
		}
	}
	return out
}

// dictMarshaler encodes a FieldKindDict into a zap object without allocations where possible.
type dictMarshaler struct{ fs []logstox.Field }

func (d dictMarshaler) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, f := range d.fs {
		if f.IsSkip() {
			continue
		}
		switch f.Kind() {
		case logstox.FieldKindString:
			enc.AddString(f.Key, f.Value.(string))
		case logstox.FieldKindBool:
			enc.AddBool(f.Key, f.Value.(bool))
		case logstox.FieldKindInt64:
			enc.AddInt64(f.Key, f.Value.(int64))
		case logstox.FieldKindUint64:
			enc.AddUint64(f.Key, f.Value.(uint64))
		case logstox.FieldKindFloat64:
			enc.AddFloat64(f.Key, f.Value.(float64))
		case logstox.FieldKindDuration:
			enc.AddDuration(f.Key, f.Value.(time.Duration))
		case logstox.FieldKindTime:
			enc.AddTime(f.Key, f.Value.(time.Time))
		case logstox.FieldKindError:
			// As a string to keep nested objects simple; top-level uses zap.NamedError/zap.Error.
			enc.AddString(f.Key, f.Value.(error).Error())
		case logstox.FieldKindStrings:
			enc.AddArray(f.Key, stringArray(f.Value.([]string)))
		case logstox.FieldKindBools:
			enc.AddArray(f.Key, boolArray(f.Value.([]bool)))
		case logstox.FieldKindInt64s:
			enc.AddArray(f.Key, int64Array(f.Value.([]int64)))
		case logstox.FieldKindUint64s:
			enc.AddArray(f.Key, uint64Array(f.Value.([]uint64)))
		case logstox.FieldKindFloat64s:
			enc.AddArray(f.Key, float64Array(f.Value.([]float64)))
		case logstox.FieldKindErrors:
			enc.AddArray(f.Key, errorArray(f.Value.([]error)))
		case logstox.FieldKindDict:
			enc.AddObject(f.Key, dictMarshaler{fs: f.Value.([]logstox.Field)})
		case logstox.FieldKindRawJSON:
			// Preserve raw JSON in nested objects.
			enc.AddReflected(f.Key, json.RawMessage(f.Value.([]byte)))
		case logstox.FieldKindHexBytes:
			enc.AddString(f.Key, hex.EncodeToString(f.Value.([]byte)))
		case logstox.FieldKindTimestamp:
			t := f.Value.(time.Time)
			if t.IsZero() {
				t = time.Now()
			}
			enc.AddTime(f.Key, t)
		case logstox.FieldKindAny:
			enc.AddReflected(f.Key, f.Value)
		case logstox.FieldKindLazyValue:
			// Lazy funcs are expanded at the call site with level checks, so ignore here.
			continue
		default:
			// skip invalid
		}
	}
	return nil
}

// Arrays for dictMarshaler
type (
	stringArray  []string
	boolArray    []bool
	int64Array   []int64
	uint64Array  []uint64
	float64Array []float64
	errorArray   []error
)

func (a stringArray) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, v := range a {
		enc.AppendString(v)
	}
	return nil
}
func (a boolArray) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, v := range a {
		enc.AppendBool(v)
	}
	return nil
}
func (a int64Array) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, v := range a {
		enc.AppendInt64(v)
	}
	return nil
}
func (a uint64Array) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, v := range a {
		enc.AppendUint64(v)
	}
	return nil
}
func (a float64Array) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, v := range a {
		enc.AppendFloat64(v)
	}
	return nil
}
func (a errorArray) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, v := range a {
		if v != nil {
			enc.AppendString(v.Error())
		} else {
			enc.AppendString("<nil>")
		}
	}
	return nil
}

// Escape hatch

// Native lets callers pass a zap.Field directly through logstox.
// Use for zap-specific features (e.g., zap.Stack, zap.Reflect).
// Example: log.Error("oops", zapx.Native(zap.Stack("stack")))
type native struct{ zf zap.Field }

func Native(zf zap.Field) logstox.Field {
	// Use FieldKindAny; mapper handles the 'native' unwrapping first.
	return logstox.Any("__zap_native__", native{zf: zf})
}
