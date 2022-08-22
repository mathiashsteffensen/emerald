package vm

import "testing"

func TestErrorHandling(t *testing.T) {
	tests := []vmTestCase{
		// Error raising not implemented
		//{
		//	name: "calling method with wrong amount of arguments",
		//	input: `
		//	def say_hello(name)
		//		puts("Hello " + name + "!")
		//	end
		//
		//	say_hello
		//	`,
		//	expected: "error:wrong number of arguments (given 0, expected 1)",
		//},
	}

	runVmTests(t, tests)
}
