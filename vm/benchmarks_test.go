package vm

import (
	"emerald/compiler"
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
	
	fib(18)
	`

	for i := 0; i < b.N; i++ {
		program := parse(input)
		comp := compiler.New()

		comp.Compile(program)

		vm := New("", comp.Bytecode())
		vm.Run()
	}
}
