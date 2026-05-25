---
name: relax
description: 'Use when working with github.com/luckyman42/relax, choosing Check*, FailOnError*, HandleFailer, errors.As with Failer, or explaining panic-based error propagation in Go trusted internal layers.'
argument-hint: 'Describe the integration point or question about using relax'
---

# Relax

This skill teaches an agent how to use `github.com/luckyman42/relax` correctly.

Use it when writing code with the library, reviewing code that already uses it, or explaining the library to someone who is new to it.

This skill is **self-contained**. It includes the full public API surface with types and signatures so that an AI agent can generate correct code without access to the source.

## Installation and Import

```bash
go get github.com/luckyman42/relax
```

```go
import "github.com/luckyman42/relax"
```

Requires **Go 1.24+**.

## When to Use

- Integrating `github.com/luckyman42/relax` into a Go codebase
- Choosing between `FailWith`, `FailOnError*`, `Check*`, and `HandleFailer`
- Explaining the mental model of panic-based propagation inside trusted internal layers
- Showing how to use `errors.As` or `errors.Is` without depending on `Failer`
- Designing goroutine boundaries with `go relax.HandleFailer(fn, onError)`
- Explaining why helper support stops at three non-error return values

## Mental Model

- Inside trusted internal layers, use `FailWith` or `FailOnError*` to stop local forwarding boilerplate.
- At the first boundary that should return a normal Go `error`, use a `Check*` helper.
- In goroutines, keep the goroutine launch explicit and use `go relax.HandleFailer(fn, onError)`.
- Only panics carrying a `Failer` are recovered.
- Non-`Failer` panics must propagate unchanged.
- `Failer` implements `Unwrap()`, so callers usually do not need to know about `Failer` to inspect the original error.

## Complete API Reference

### Failer Type

```go
type Failer struct {
	Err     error
	Context map[string]any
}
```

**Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `Err` | `error` | The original wrapped error |
| `Context` | `map[string]any` | Optional structured key/value metadata |

**Methods:**

```go
// Error returns f.Err.Error(). Returns "" if Err is nil.
func (f Failer) Error() string

// Unwrap returns f.Err — enables errors.Is and errors.As on the inner error.
func (f Failer) Unwrap() error

// Fail re-panics this Failer. Optional keyVals are merged into Context.
// This avoids wrapping a Failer inside another Failer when rethrowing.
func (f Failer) Fail(keyVals ...any)
```

`Failer` satisfies the `error` interface. Because it implements `Unwrap()`, standard `errors.Is` and `errors.As` calls transparently match the inner `Err` without the caller needing to know about `Failer`.

### Failure Propagation Functions

```go
// FailWith panics with a Failer wrapping err.
// If err is nil, FailWith is a no-op (returns immediately).
// If err is already a Failer (value or pointer), it is re-panicked directly;
// any provided keyVals are merged into the existing Context.
// keyVals are alternating key, value pairs. Keys are stringified via fmt.Sprint.
// An odd number of keyVals is allowed; the final key gets a nil value.
func FailWith(err error, keyVals ...any)

// FailOnError returns v if err == nil; otherwise calls FailWith(err).
func FailOnError[T any](v T, err error) T

// FailOnError2 returns v1, v2 if err == nil; otherwise calls FailWith(err).
func FailOnError2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2)

// FailOnError3 returns v1, v2, v3 if err == nil; otherwise calls FailWith(err).
func FailOnError3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3)
```

### Recovery Boundary Functions

```go
// CheckFailer executes fn; recovers Failer panics into error. Other panics re-panic.
func CheckFailer(fn func()) error

// CheckError executes fn which returns error. Recovers Failer panics into error.
// If fn returns a non-nil error normally, that error is returned unchanged.
func CheckError(fn func() error) error

// CheckValue executes fn returning T. Recovers Failer panics into error.
func CheckValue[T any](fn func() T) (T, error)

// CheckValue2 executes fn returning (T1, T2). Recovers Failer panics into error.
func CheckValue2[T1, T2 any](fn func() (T1, T2)) (T1, T2, error)

// CheckValue3 executes fn returning (T1, T2, T3). Recovers Failer panics into error.
func CheckValue3[T1, T2, T3 any](fn func() (T1, T2, T3)) (T1, T2, T3, error)

// CheckResult executes fn returning (T, error). Recovers Failer panics into error.
// If fn returns a non-nil error normally, that error is returned unchanged.
func CheckResult[T any](fn func() (T, error)) (T, error)

// CheckResult2 executes fn returning (T1, T2, error). Recovers Failer panics into error.
func CheckResult2[T1, T2 any](fn func() (T1, T2, error)) (T1, T2, error)

// CheckResult3 executes fn returning (T1, T2, T3, error). Recovers Failer panics into error.
func CheckResult3[T1, T2, T3 any](fn func() (T1, T2, T3, error)) (T1, T2, T3, error)

// HandleFailer executes fn inside a CheckFailer boundary and forwards any
// recovered Failer as a standard error to onError.
// Only Failer panics are recovered; other panics re-panic.
// onError must not be nil — passing nil causes an immediate panic.
func HandleFailer(fn func(), onError func(error))
```

### Utility Functions

```go
// ConvertToFailer converts any error into a Failer value.
// If err is already a Failer (or wraps one via errors.As), returns it unchanged.
// Otherwise creates a new Failer wrapping err.
// If err is nil, returns a zero-value Failer.
func ConvertToFailer(err error) Failer

// IsFailer reports whether err is or wraps a Failer.
// Returns false if err is nil.
func IsFailer(err error) bool
```

## Helper Selection

Use the helper that matches the function shape.

- `func()` -> `CheckFailer`
- `func() error` -> `CheckError`
- `func() T` -> `CheckValue`
- `func() (T1, T2)` -> `CheckValue2`
- `func() (T1, T2, T3)` -> `CheckValue3`
- `func() (T, error)` -> `CheckResult`
- `func() (T1, T2, error)` -> `CheckResult2`
- `func() (T1, T2, T3, error)` -> `CheckResult3`
- `(T, error)` -> `FailOnError`
- `(T1, T2, error)` -> `FailOnError2`
- `(T1, T2, T3, error)` -> `FailOnError3`

Support intentionally stops at three non-error return values because Go does not provide a generic abstraction for arbitrary return arity. For `4+` values, wrap the values in a struct or dedicated result type and return fewer parameters.

## Transparent Error Compatibility (Key Insight)

The most important thing to understand: **callers at the boundary do not need to change their error-handling logic at all.**

When a `Check*` function recovers a `Failer`, it returns it as a standard `error`. Because `Failer` implements `Unwrap()`, all existing `errors.Is` and `errors.As` calls against domain error types **continue to work unchanged**.

```go
// Before relax — explicit error forwarding:
func ProcessOrder(id int) error {
	order, err := loadOrder(id)
	if err != nil {
		return err
	}
	return validateOrder(order)
}

// After relax — same caller-side handling, zero changes needed:
func ProcessOrder(id int) error {
	return relax.CheckError(func() error {
		order := relax.FailOnError(loadOrder(id))
		return validateOrder(order)
	})
}

// The caller of ProcessOrder does NOT change at all:
err := ProcessOrder(42)
var notFound *OrderNotFoundError
if errors.As(err, &notFound) {
	// works exactly as before — Failer is invisible to this code
	http.Error(w, notFound.Error(), 404)
}
```

**This means adopting `relax` is minimally invasive:**

- The boundary function's **signature** stays the same (`error` return).
- The **caller's** `errors.As` / `errors.Is` logic is **unchanged**.
- The caller never needs to import `relax` or know about `Failer`.
- Only the **internal implementation** switches from explicit `if err != nil { return err }` chains to `FailOnError*` / `FailWith`.

When generating code that adopts `relax`, **do not modify the caller-side error handling**. Only refactor the internal layers.

## Core Patterns

### 1. Boundary Pattern

Use `FailOnError*` or `FailWith` internally, then recover once at the outer boundary.

```go
func HandleRequest(id int) error {
	return relax.CheckError(func() error {
		user := relax.FailOnError(loadUser(id))
		processUser(user)
		return nil
	})
}
```

### 2. Fail Fast Inside Trusted Layers

If an intermediate layer has no meaningful recovery decision, it should not forward the same error explicitly.

```go
func A() error { return relax.CheckFailer(B) }
func B()       { C() }
func C()       { D() }
func D()       { relax.FailWith(E()) }
func E() error { return errors.New("storage unavailable") }
```

### 3. Work With Underlying Errors

Prefer matching the original domain error with `errors.As` or `errors.Is`.

```go
type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return "invalid " + e.Field
}

_, err := relax.CheckValue(func() string {
	relax.FailWith(&ValidationError{Field: "email"})
	return ""
})

var target *ValidationError
if errors.As(err, &target) {
	fmt.Println(target.Field)
}
```

Only inspect `Failer` directly when you need `Context`. For everything else, use `errors.As` or `errors.Is` against the inner domain error.

### 4. Add Structured Context

```go
if err := saveUser(user); err != nil {
	relax.FailWith(err,
		"user_id", user.ID,
		"operation", "save_user",
	)
}
```

If the error is already a `Failer`, context is merged instead of double-wrapping.

Keys are stringified using `fmt.Sprint`; the last write wins for a given key. To avoid collisions when multiple layers annotate the same failure, use namespaced keys:

```go
relax.FailWith(err, "db.operation", "save")
relax.FailWith(err, "http.operation", "POST /users")
```

### 5. Goroutine Boundary

```go
go relax.HandleFailer(func() {
	user := relax.FailOnError(loadUser(id))
	syncUser(user)
}, func(err error) {
	log.Printf("worker failed: %v", err)
})
```

Use `HandleFailer` directly when the code runs synchronously and you want an explicit local handler.

## Error Integration

Three patterns for using `relax` with standard Go error types:

**`fmt.Errorf` with `%w`** — add context while keeping the error chain intact:

```go
relax.FailWith(fmt.Errorf("loading user %d: %w", id, err))
// errors.As reaches the original error through Failer.Unwrap() and fmt.Errorf's chain.
```

**Sentinel errors with `errors.Is`** — works automatically through `Failer.Unwrap()`:

```go
var ErrNotFound = errors.New("not found")

relax.FailWith(ErrNotFound)

// At the boundary — true because Failer.Unwrap() returns ErrNotFound:
errors.Is(err, ErrNotFound)
```

**`errors.Join` (Go 1.20+)** — when multiple errors occur inside a boundary:

```go
relax.FailWith(errors.Join(errValidation, errDatabase))
// Both errors are reachable via errors.As through the Unwrap() []error interface.
```

If stack traces are needed, wrap the error before passing it to `FailWith`:

```go
relax.FailWith(fmt.Errorf("operation failed: %w", err))
```

The wrapped error remains accessible via `errors.As` and `errors.Is`. `relax` does not endorse any specific stack trace library.

## Guidance for Agents

When generating code or advice for this library:

1. Identify the boundary where normal Go `error` values should reappear.
2. Keep `FailWith` and `FailOnError*` inside trusted internal layers only.
3. Prefer `errors.As` and `errors.Is` against the inner domain error before recommending direct `Failer` inspection.
4. Mention the `3`-value arity cutoff and recommend struct wrapping for `4+` returns.
5. Keep it explicit that this library is not a replacement for Go's normal public API error handling.

## Anti-Patterns

- Do not recommend `relax` for exported public APIs by default.
- Do not swallow non-`Failer` panics.
- Do not tell users they must depend on `Failer` if they only need the original error.
- Do not suggest adding `CheckValue4`, `CheckResult4`, and so on unless the user explicitly wants to expand the public API.
- Do not use this library as ordinary control flow.

## Good Fit

- service-layer orchestration
- request and command pipelines
- background jobs and workers
- CLI flows
- deep internal call chains where middle layers would only re-return the same error

## Poor Fit

- exported public APIs
- low-level reusable libraries intended for general consumption
- hot performance-critical loops
- code paths where explicit `error` handling is clearer than propagation

## Design Guarantees

The library guarantees the following invariants:

1. Only `Failer` panics are recovered by `Check*` and `HandleFailer`.
2. Non-`Failer` panics (programmer bugs, runtime faults) propagate unchanged.
3. The original error is always preserved through `Failer.Unwrap()`.
4. Existing `Failer` values are never double-wrapped — `FailWith` on a `Failer` re-panics it directly.
5. `errors.Is` and `errors.As` work transparently against the inner error.
6. `FailWith(nil)` is always a no-op.
7. `HandleFailer` with a nil `onError` panics immediately (fail-fast on misconfiguration).

## Behavioral Details

### FailWith Semantics

- `FailWith(nil)` → no-op, returns immediately.
- `FailWith(normalErr)` → creates a new `Failer{Err: normalErr}` and panics with it.
- `FailWith(existingFailer)` → calls `existingFailer.Fail()`, re-panicking the same `Failer` (no new wrapping).
- `FailWith(err, "key1", val1, "key2", val2)` → attaches key/value pairs to `Failer.Context`.

### Failer.Fail Semantics

- Merges any provided key/value pairs into the existing `Context` map.
- Then panics with the (possibly enriched) `Failer` value.
- Used internally by `FailWith` when the error is already a `Failer`.

### Recovery Mechanics

All `Check*` functions use `defer recover()` internally. If the recovered value is a `Failer` (value or pointer), it is returned as the `error`. If the recovered value is anything else (including nil), it is re-panicked. This means:

- A `Check*` boundary **never** silently swallows a non-`Failer` panic.
- A `Check*` boundary **always** converts a `Failer` into a normal `error` return.

### CheckError vs CheckFailer

- `CheckFailer(fn)` — use when `fn` has no return value. Any failure comes only from `FailWith`/`FailOnError*`.
- `CheckError(fn)` — use when `fn` already returns `error`. Both the returned error and a recovered `Failer` panic are merged into the single error return.

### CheckValue vs CheckResult

- `CheckValue[T](fn)` — use when `fn` returns only value(s), no error. Failures come only from panics.
- `CheckResult[T](fn)` — use when `fn` returns `(T, error)`. Both the returned error and a recovered `Failer` are merged into the error return.

## Flow Diagram

```
Caller → A (boundary) → B → C → D → E

E returns error
D calls FailWith(err) → panic(Failer)
C does nothing (panic unwinds through it)
B does nothing (panic unwinds through it)
A uses Check* → recovers Failer → returns error to Caller
```

The key insight: B and C do not participate in error forwarding at all.