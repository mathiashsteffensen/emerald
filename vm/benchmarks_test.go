package vm

import (
	"emerald/compiler"
	"emerald/core"
	"testing"
)

func BenchmarkFibonacci(b *testing.B) {
	input := `
	def fib(n)
	  if n == 0
		return n
	  end
	  if  n == 1
		return n
	  end
	
	  fib(n - 1) + fib(n - 2)
	end
	
	fib(28)
	`

	for i := 0; i < b.N; i++ {
		program := parse(input)
		comp := compiler.New()

		err := comp.Compile(program)
		if err != nil {
			b.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())

		err = vm.Run()
		if err != nil {
			b.Fatalf("vm error: %s", err)
		}

		val := vm.LastPoppedStackElem().(*core.IntegerInstance).Value
		if val != 317811 {
			b.Fatalf("Fibonacci calculation wrong, want=317811 got=%d", val)
		}
	}
}
