package emerald

import (
	"testing"
)

/*
Baseline benchmarks for integer arithmetic.
Run with: go test -bench=BenchmarkIntegerArithmetic -benchmem
*/
func BenchmarkIntegerArithmetic(b *testing.B) {
	input := `
	i = 0
	while i < 1000
		i = i + 1
	end
	`

	engine := New()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.Eval(input)
	}
}
