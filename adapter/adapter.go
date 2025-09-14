package adapter

import (
	"github.com/khinshankhan/logstox"
)

// Adapter bridges a Logger[Base] so callers can log with fields of type App.
type Adapter[Base, App any] struct {
	// Base is the underlying logger we ultimately call.
	Base logstox.Logger[Base]
	// ToBase converts an application field into the Base field.
	ToBase func(App) Base
}

// mapSlice applies f to all items of in, returning a freshly allocated slice.
// A nil or empty input yields nil for zero allocations where downstream APIs
// treat nil and empty equally.
func mapSlice[Out, In any](f func(In) Out, in []In) []Out {
	if len(in) == 0 {
		return nil
	}
	out := make([]Out, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

// DEBUG (-1): for recording messages useful for debugging.
func (a Adapter[Base, App]) Debug(msg string, fields ...App) {
	a.Base.Debug(msg, mapSlice(a.ToBase, fields)...)
}

// INFO (0): for messages describing normal application operations.
func (a Adapter[Base, App]) Info(msg string, fields ...App) {
	a.Base.Info(msg, mapSlice(a.ToBase, fields)...)
}

// WARN (1): for recording messages indicating something unusual happened that may need attention before it escalates to a more severe issue.
func (a Adapter[Base, App]) Warn(msg string, fields ...App) {
	a.Base.Warn(msg, mapSlice(a.ToBase, fields)...)
}

// ERROR (2): for recording unexpected error conditions in the program.
func (a Adapter[Base, App]) Error(msg string, fields ...App) {
	a.Base.Error(msg, mapSlice(a.ToBase, fields)...)
}

// DPANIC (3): for recording severe error conditions in development. It behaves like PANIC in development and ERROR in production.
func (a Adapter[Base, App]) DPanic(msg string, fields ...App) {
	a.Base.DPanic(msg, mapSlice(a.ToBase, fields)...)
}

// PANIC (4): calls panic() after logging an error condition.
func (a Adapter[Base, App]) Panic(msg string, fields ...App) {
	a.Base.Panic(msg, mapSlice(a.ToBase, fields)...)
}

// FATAL (5): calls os.Exit(1) after logging an error condition.
func (a Adapter[Base, App]) Fatal(msg string, fields ...App) {
	a.Base.Fatal(msg, mapSlice(a.ToBase, fields)...)
}

// With returns a new Adapter with Base.With(...) applied and the same converter.
func (a Adapter[Base, App]) With(fields ...App) logstox.Logger[App] {
	return Adapter[Base, App]{
		Base:   a.Base.With(mapSlice(a.ToBase, fields)...),
		ToBase: a.ToBase,
	}
}

// Named returns a new Adapter with Base.Named(name) and the same converter.
func (a Adapter[Base, App]) Named(name string) logstox.Logger[App] {
	return Adapter[Base, App]{
		Base:   a.Base.Named(name),
		ToBase: a.ToBase,
	}
}

// Sync delegates to the underlying logger's Sync.
func (a Adapter[Base, App]) Sync() error {
	return a.Base.Sync()
}
