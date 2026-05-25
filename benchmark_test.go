package relax_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/luckyman42/relax"
)

var sink error
var sinkInt int
var errFixed = errors.New("benchmark error")

func explicitDeep_ok(n int) (int, error) {
	if n == 0 {
		return 42, nil
	}
	v, err := explicitDeep_ok(n - 1)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func relaxDeep_ok(n int) int {
	if n == 0 {
		return 42
	}
	return relaxDeep_ok(n - 1)
}

func explicitDeep_err(n int) (int, error) {
	if n == 0 {
		return 0, errFixed
	}
	v, err := explicitDeep_err(n - 1)
	if err != nil {
		return 0, err
	}
	return v, nil
}

func relaxDeep_err(n int) {
	if n == 0 {
		relax.FailWith(errFixed)
		return
	}
	relaxDeep_err(n - 1)
}

// BenchmarkHappyPath shows that relax overhead is one defer at the boundary.
// Explicit accumulates (T, error) return overhead per frame; relax does not.
// The two cross at depth ≈ 8.
func BenchmarkHappyPath(b *testing.B) {
	for _, depth := range []int{1, 5, 8, 10} {
		b.Run(fmt.Sprintf("explicit/depth=%d", depth), func(b *testing.B) {
			for b.Loop() {
				sinkInt, sink = explicitDeep_ok(depth)
			}
		})
		b.Run(fmt.Sprintf("relax/depth=%d", depth), func(b *testing.B) {
			for b.Loop() {
				sinkInt, sink = relax.CheckValue(func() int { return relaxDeep_ok(depth) })
			}
		})
	}
}

// BenchmarkErrorPath shows that panic+recover is a constant ~400-700 ns cost
// regardless of depth — the mechanism dominates, not the stack unwind.
func BenchmarkErrorPath(b *testing.B) {
	for _, depth := range []int{1, 5, 10} {
		b.Run(fmt.Sprintf("explicit/depth=%d", depth), func(b *testing.B) {
			for b.Loop() {
				sinkInt, sink = explicitDeep_err(depth)
			}
		})
		b.Run(fmt.Sprintf("relax/depth=%d", depth), func(b *testing.B) {
			for b.Loop() {
				sink = relax.CheckFailer(func() { relaxDeep_err(depth) })
			}
		})
	}
}

// BenchmarkWithContext measures the additional cost of attaching structured
// metadata to a failure via FailWith key-value pairs.
func BenchmarkWithContext(b *testing.B) {
	b.Run("no_context", func(b *testing.B) {
		for b.Loop() {
			sink = relax.CheckFailer(func() { relax.FailWith(errFixed) })
		}
	})
	b.Run("4_key_value_pairs", func(b *testing.B) {
		for b.Loop() {
			sink = relax.CheckFailer(func() {
				relax.FailWith(errFixed, "op", "save", "id", 42, "table", "users", "db", "primary")
			})
		}
	})
}
