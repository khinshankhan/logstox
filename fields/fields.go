package fields

import (
	"context"
	"time"
)

// FieldKind describes the concrete type carried by a Field's value.
type FieldKind uint8

const (
	FieldKindInvalid FieldKind = iota // zero / no-op
	FieldKindAny

	// Scalars
	FieldKindString
	FieldKindBool
	FieldKindInt64
	FieldKindUint64
	FieldKindFloat64
	FieldKindTime
	FieldKindDuration
	FieldKindError

	// Slices
	FieldKindStrings
	FieldKindBools
	FieldKindInt64s
	FieldKindUint64s
	FieldKindFloat64s
	FieldKindErrors

	// Special
	FieldKindDict       // sub-fields (Value is []Field)
	FieldKindRawJSON    // []byte that is already JSON
	FieldKindHexBytes   // []byte to render as hex string
	FieldKindLazyFields // lazy: func(context.Context) []Field
	FieldKindLazyValue  // lazy: func() []Field
	FieldKindTimestamp  // backend inserts current timestamp (or uses Value as time.Time if provided)
)

// Conventional keys used by helpers.
const (
	ErrorKey     = "error"
	TimestampKey = "ts"
)

// Field is a portable structured field: a key plus a typed value.
// The unexported 'kind' enforces invariants via the constructors below.
type Field struct {
	Key   string
	kind  FieldKind
	Value any
}

// Kind returns the field's discriminant for backend mapping.
func (f Field) Kind() FieldKind {
	return f.kind
}

func Nop() Field {
	return Field{
		Key:   "",
		kind:  FieldKindInvalid,
		Value: struct{}{},
	}
}

// IsZero reports whether f is a no-op field.
func (f Field) IsZero() bool {
	return f.kind == FieldKindInvalid
}

// IsSkip reports whether the field should be emitted.
func (f Field) IsSkip() bool {
	return !f.IsZero()
}

// Scalars

func Any(k string, v any) Field {
	return Field{Key: k, kind: FieldKindAny, Value: v}
}
func String(k, v string) Field {
	return Field{Key: k, kind: FieldKindString, Value: v}
}
func Bool(k string, v bool) Field {
	return Field{Key: k, kind: FieldKindBool, Value: v}
}
func Int(k string, v int) Field {
	return Field{Key: k, kind: FieldKindInt64, Value: int64(v)}
}
func Int64(k string, v int64) Field {
	return Field{Key: k, kind: FieldKindInt64, Value: v}
}
func Uint(k string, v uint) Field {
	return Field{Key: k, kind: FieldKindUint64, Value: uint64(v)}
}
func Uint64(k string, v uint64) Field {
	return Field{Key: k, kind: FieldKindUint64, Value: v}
}
func Float64(k string, v float64) Field {
	return Field{Key: k, kind: FieldKindFloat64, Value: v}
}
func TimeField(k string, v time.Time) Field {
	return Field{Key: k, kind: FieldKindTime, Value: v}
}
func Duration(k string, v time.Duration) Field {
	return Field{Key: k, kind: FieldKindDuration, Value: v}
}

// Errors

// Error adds a non-nil error under the conventional key ("error").
// If err is nil, returns a no-op.
func Error(err error) Field {
	if err == nil {
		return Nop()
	}
	return Field{Key: ErrorKey, kind: FieldKindError, Value: err}
}

// NamedError adds a non-nil error with a custom key.
// If err is nil, returns a no-op.
func NamedError(k string, err error) Field {
	if err == nil {
		return Nop()
	}
	return Field{Key: k, kind: FieldKindError, Value: err}
}

// Slices (not copied; pass a copy if you may mutate later)

func Strings(k string, v []string) Field {
	return Field{Key: k, kind: FieldKindStrings, Value: v}
}
func Bools(k string, v []bool) Field {
	return Field{Key: k, kind: FieldKindBools, Value: v}
}
func Int64s(k string, v []int64) Field {
	return Field{Key: k, kind: FieldKindInt64s, Value: v}
}
func Uint64s(k string, v []uint64) Field {
	return Field{Key: k, kind: FieldKindUint64s, Value: v}
}
func Float64s(k string, v []float64) Field {
	return Field{Key: k, kind: FieldKindFloat64s, Value: v}
}
func Errors(k string, v []error) Field {
	return Field{Key: k, kind: FieldKindErrors, Value: v}
}

// Special

// Dict groups sub-fields under a single key (zap: Object, slog: Group).
// NOTE: this does not copy the slice; pass a copy if you will mutate it.
func Dict(k string, fields ...Field) Field {
	return Field{Key: k, kind: FieldKindDict, Value: fields}
}

// RawJSON inserts pre-encoded JSON bytes under key.
// NOTE: Backends that don’t support raw JSON may encode it as a string or bytes (eg zap supports; slog may treat as
// []byte/string).
func RawJSON(k string, json []byte) Field {
	return Field{Key: k, kind: FieldKindRawJSON, Value: json}
}

// Hex encodes []byte as a lowercase hexadecimal string at the backend.
func Hex(k string, b []byte) Field {
	return Field{Key: k, kind: FieldKindHexBytes, Value: b}
}

// LazyFields runs lazily at log time (only if enabled) and returns extra fields to append.
// NOTE: The function should be fast and side-effect free.
func LazyFields(fn func(context.Context) []Field) Field {
	if fn == nil {
		return Field{}
	}
	return Field{kind: FieldKindLazyFields, Value: fn}
}

// Convenience for no-ctx callers.
// NOTE: The function should be fast and side-effect free.
func Lazy(fn func() []Field) Field {
	if fn == nil {
		return Field{}
	}
	return LazyFields(func(context.Context) []Field { return fn() })
}

// Timestamp asks the backend to attach a timestamp field.
// If t is zero, backends should use time.Now(); otherwise use t.
// The key defaults to TimestampKey ("ts"); backends may honor Options.TimeLayout.
func Timestamp(t time.Time) Field {
	return Field{Key: TimestampKey, kind: FieldKindTimestamp, Value: t}
}

// TimestampAt is the same as Timestamp but with a custom key.
func TimestampAt(k string, t time.Time) Field {
	return Field{Key: k, kind: FieldKindTimestamp, Value: t}
}

// From chooses a FieldKind for common types; otherwise returns Any.
func From(k string, v any) Field {
	switch t := v.(type) {
	case string:
		return String(k, t)
	case bool:
		return Bool(k, t)
	case int:
		return Int(k, t)
	case int64:
		return Int64(k, t)
	case uint:
		return Uint(k, t)
	case uint64:
		return Uint64(k, t)
	case float64:
		return Float64(k, t)
	case time.Time:
		return TimeField(k, t)
	case time.Duration:
		return Duration(k, t)
	case error:
		return NamedError(k, t) // respects custom key
	default:
		return Any(k, v)
	}
}
