package zapx

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/khinshankhan/logstox/fields"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ToZap(f fields.Field) zap.Field {
	switch f.Kind() {
	case fields.FieldKindString:
		return zap.String(f.Key, f.Value.(string))
	case fields.FieldKindBool:
		return zap.Bool(f.Key, f.Value.(bool))
	case fields.FieldKindInt64:
		return zap.Int64(f.Key, f.Value.(int64))
	case fields.FieldKindUint64:
		return zap.Uint64(f.Key, f.Value.(uint64))
	case fields.FieldKindFloat64:
		return zap.Float64(f.Key, f.Value.(float64))
	case fields.FieldKindDuration:
		return zap.Duration(f.Key, f.Value.(time.Duration))
	case fields.FieldKindTime:
		return zap.Time(f.Key, f.Value.(time.Time))
	case fields.FieldKindError:
		if f.Key == "" || f.Key == fields.ErrorKey {
			return zap.Error(f.Value.(error))
		}
		return zap.NamedError(f.Key, f.Value.(error))
	case fields.FieldKindStrings:
		return zap.Strings(f.Key, f.Value.([]string))
	case fields.FieldKindBools:
		return zap.Bools(f.Key, f.Value.([]bool))
	case fields.FieldKindInt64s:
		return zap.Int64s(f.Key, f.Value.([]int64))
	case fields.FieldKindUint64s:
		return zap.Uint64s(f.Key, f.Value.([]uint64))
	case fields.FieldKindFloat64s:
		return zap.Float64s(f.Key, f.Value.([]float64))
	case fields.FieldKindErrors:
		return zap.Errors(f.Key, f.Value.([]error))
	case fields.FieldKindRawJSON:
		return zap.Any(f.Key, json.RawMessage(f.Value.([]byte)))
	case fields.FieldKindHexBytes:
		return zap.String(f.Key, hex.EncodeToString(f.Value.([]byte)))
	case fields.FieldKindDict:
		return zap.Object(f.Key, dict{f.Value.([]fields.Field)})
	case fields.FieldKindTimestamp:
		t := f.Value.(time.Time)
		if t.IsZero() {
			t = time.Now()
		}
		return zap.Time(f.Key, t)
	case fields.FieldKindAny:
		return zap.Any(f.Key, f.Value)
	default:
		// TODO: look into exhaustive checks
		return zap.Skip()
	}
}

type dict struct{ fs []fields.Field }

func (d dict) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, f := range d.fs {
		switch f.Kind() {
		case fields.FieldKindString:
			enc.AddString(f.Key, f.Value.(string))
		case fields.FieldKindBool:
			enc.AddBool(f.Key, f.Value.(bool))
		case fields.FieldKindInt64:
			enc.AddInt64(f.Key, f.Value.(int64))
		case fields.FieldKindUint64:
			enc.AddUint64(f.Key, f.Value.(uint64))
		case fields.FieldKindFloat64:
			enc.AddFloat64(f.Key, f.Value.(float64))
		case fields.FieldKindDuration:
			enc.AddDuration(f.Key, f.Value.(time.Duration))
		case fields.FieldKindTime:
			enc.AddTime(f.Key, f.Value.(time.Time))
		case fields.FieldKindError:
			enc.AddString(f.Key, f.Value.(error).Error())
		case fields.FieldKindStrings:
			enc.AddArray(f.Key, stringArray(f.Value.([]string)))
		case fields.FieldKindBools:
			enc.AddArray(f.Key, boolArray(f.Value.([]bool)))
		case fields.FieldKindInt64s:
			enc.AddArray(f.Key, int64Array(f.Value.([]int64)))
		case fields.FieldKindUint64s:
			enc.AddArray(f.Key, uint64Array(f.Value.([]uint64)))
		case fields.FieldKindFloat64s:
			enc.AddArray(f.Key, float64Array(f.Value.([]float64)))
		case fields.FieldKindErrors:
			enc.AddArray(f.Key, errorArray(f.Value.([]error)))
		case fields.FieldKindDict:
			enc.AddObject(f.Key, dict{f.Value.([]fields.Field)})
		case fields.FieldKindRawJSON:
			enc.AddReflected(f.Key, json.RawMessage(f.Value.([]byte)))
		case fields.FieldKindHexBytes:
			enc.AddString(f.Key, hex.EncodeToString(f.Value.([]byte)))
		case fields.FieldKindTimestamp:
			t := f.Value.(time.Time)
			if t.IsZero() {
				t = time.Now()
			}
			enc.AddTime(f.Key, t)
		case fields.FieldKindAny:
			enc.AddReflected(f.Key, f.Value)
		}
	}
	return nil
}

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
