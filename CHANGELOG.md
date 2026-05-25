# Changelog

All notable changes to this project will be documented in this file.

## [1.0.0] - 2026-05-25

First stable release. The public API is now considered stable.

### Added

- `FailWith` — panics with a `Failer` wrapping any `error`; merges context if the error is already a `Failer`
- `FailOnError`, `FailOnError2`, `FailOnError3` — inline helpers for the common `(T, error)` return shapes
- `CheckFailer`, `CheckError` — recovery boundaries for `func()` and `func() error`
- `CheckValue`, `CheckValue2`, `CheckValue3` — recovery boundaries for value-returning closures
- `CheckResult`, `CheckResult2`, `CheckResult3` — recovery boundaries for closures that already return `(T, error)`
- `HandleFailer` — explicit error handler boundary; designed for goroutine entry points and worker loops
- `Failer` type with `Err` (preserved original error) and `Context` (structured key/value metadata)
- `Failer.Fail` — re-panics an existing `Failer`, merging additional context without double-wrapping
- `ConvertToFailer`, `IsFailer` — utility helpers for inspecting `Failer` values
- Structured context support via `FailWith(err, "key", value, ...)` with namespaced key convention
- `Failer` implements `Unwrap()` so `errors.Is` and `errors.As` work transparently against the inner error
- Reusable AI skill in `.github/skills/relax/` with the full API surface
- Benchmark suite comparing relax vs explicit propagation at depths 1, 5, 8, 10

### Design guarantees

- Only `Failer` panics are recovered; non-`Failer` panics propagate unchanged
- `FailWith(nil)` is always a no-op
- Existing `Failer` values are never double-wrapped
- `HandleFailer` with a nil `onError` panics immediately (fail-fast on misconfiguration)

## [0.5.1] - 2026-05-18

- Add reusable AI skill in `.github/skills/relax/`

## [0.5.0] - 2026-05-18

- Redesign the public API: establish the final naming convention (`FailWith`, `FailOnError*`, `Check*`, `HandleFailer`)

## [0.4.0] - 2026-05-17

- New API design
- Add goroutine support (`HandleFailer`)

## [0.3.1] - 2026-05-15

- Bug fixes

## [0.3.0] - 2026-05-15

- Add `GuardErr` feature (predecessor to `FailOnError`)

## [0.2.0] - 2026-05-14

- Refactoring and naming improvements

## [0.1.0] - 2026-05-14

- Initial release
