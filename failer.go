package relax

import (
	"errors"
	"fmt"
)

// Failer is the public, exported representation of a thrown failure.
//
// A `Failer` preserves the original `error` (Err) and an optional
// map[string]any Context for arbitrary key/value metadata. The library uses
// `Failer` values to implement structured panic-based propagation inside
// trusted internal call chains: callers may `panic` a `Failer` (via
// `FailWith` or `Failer.Fail`) and a `Check*` boundary will convert that panic
// back into a returned `error`.
type Failer struct {
	Err     error
	Context map[string]any
}

// Error returns the underlying error message for this Failer.
func (f Failer) Error() string {
	if f.Err != nil {
		return f.Err.Error()
	}
	return ""
}

// Unwrap returns the underlying error for compatibility with errors.As and errors.Is.
func (f Failer) Unwrap() error {
	return f.Err
}

// Fail throw this Failer.
// If extra key/value pairs are provided, they are merged into the Failer context.
// This avoids wrapping a Failer inside another Failer when rethrowing.
func (f Failer) Fail(keyVals ...any) {
	if len(keyVals) > 0 {
		if f.Context == nil {
			f.Context = make(map[string]any)
		}
		for i := 0; i < len(keyVals); i += 2 {
			key := fmt.Sprint(keyVals[i])
			var value any
			if i+1 < len(keyVals) {
				value = keyVals[i+1]
			}
			f.Context[key] = value
		}
	}
	panic(f)
}

func newFailer(err error, keyVals ...any) Failer {
	if len(keyVals) == 0 {
		return Failer{Err: err}
	}

	context := make(map[string]any, (len(keyVals)+1)/2)
	for i := 0; i < len(keyVals); i += 2 {
		key := fmt.Sprint(keyVals[i])
		var value any
		if i+1 < len(keyVals) {
			value = keyVals[i+1]
		}
		context[key] = value
	}
	return Failer{Err: err, Context: context}
}

func recoverFailer(r any) (Failer, bool) {
	switch f := r.(type) {
	case Failer:
		return f, true
	case *Failer:
		if f == nil {
			return Failer{}, false
		}
		return *f, true
	default:
		return Failer{}, false
	}
}

// FailWith panics with a `Failer` that wraps `err`.
//
// If `err` is already a `Failer` (value or pointer) it will be re-panicked
// directly; in that case any provided key/value pairs are merged into the
// existing `Failer.Context`. If `err` is nil, `FailWith` is a no-op.
//
// The `keyVals` are interpreted as alternating key, value pairs. Keys are
// stringified using `fmt.Sprint`; an odd number of `keyVals` is allowed and
// the final key will be assigned a `nil` value.
func FailWith(err error, keyVals ...any) {
	if err == nil {
		return
	}

	switch f := err.(type) {
	case Failer:
		f.Fail(keyVals...)
	case *Failer:
		if f == nil {
			return
		}
		(*f).Fail(keyVals...)
	default:
		panic(newFailer(err, keyVals...))
	}
}

// FailOnError returns `v` if `err == nil`; otherwise it throws the error via
// `FailWith(err)`.
//
// This reduces error-forwarding boilerplate inside internal call chains where
// panic-based propagation is acceptable. Prefer explicit returns in public APIs.
func FailOnError[T any](v T, err error) T {
	if err != nil {
		FailWith(err)
	}
	return v
}

// FailOnError2 returns v1 and v2 if err == nil; otherwise it throws the error via FailWith(err).
func FailOnError2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		FailWith(err)
	}
	return v1, v2
}

// FailOnError3 returns v1, v2 and v3 if err == nil; otherwise it throws the error via FailWith(err).
func FailOnError3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	if err != nil {
		FailWith(err)
	}
	return v1, v2, v3
}

// ConvertToFailer converts any error into a `Failer` value.
//
// If `err` is already a `Failer` (or wraps one), the underlying `Failer` is
// returned unchanged. Otherwise a new `Failer` wrapping `err` is returned.
func ConvertToFailer(err error) Failer {
	if err == nil {
		return Failer{}
	}

	switch f := err.(type) {
	case Failer:
		return f
	case *Failer:
		if f == nil {
			return Failer{}
		}
		return *f
	}

	var pointerFailer *Failer
	if errors.As(err, &pointerFailer) {
		if pointerFailer == nil {
			return Failer{}
		}
		return *pointerFailer
	}

	var failer Failer
	if errors.As(err, &failer) {
		return failer
	}

	return newFailer(err)
}

// IsFailer reports whether `err` is a `Failer` value or wraps a `Failer`.
func IsFailer(err error) bool {
	if err == nil {
		return false
	}

	var pointerFailer *Failer
	if errors.As(err, &pointerFailer) {
		return pointerFailer != nil
	}

	var failer Failer
	return errors.As(err, &failer)
}
