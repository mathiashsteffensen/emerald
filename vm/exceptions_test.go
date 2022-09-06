package vm

import "testing"

func TestErrorHandling(t *testing.T) {
	tests := []vmTestCase{
		//{
		//	name: "calling method with wrong amount of arguments",
		//	input: `
		//	def say_hello(name)
		//		puts("Hello " + name + "!")
		//	end
		//
		//	say_hello
		//	`,
		//	expected: "error:ArgumentError:wrong number of arguments (given 0, expected 1)",
		//},
		//{
		//	name: "Rescuing an error",
		//	input: `
		//		def say_hello(name)
		//			puts("Hello " + name + "!")
		//		end
		//
		//		$suppressed = 0
		//
		//		def suppress_argument_errors
		//			say_hello
		//		rescue ArgumentError
		//			$suppressed = $suppressed + 1
		//		end
		//
		//		suppress_argument_errors
		//
		//		$suppressed
		//	`,
		//	expected: 1,
		//},
		//{
		//	name: "Rescuing an error by specifying subclass",
		//	input: `
		//		def say_hello(name)
		//			puts("Hello " + name + "!")
		//		end
		//
		//		$suppressed = 0
		//
		//		def suppress_errors
		//			say_hello
		//		rescue StandardError
		//			$suppressed = $suppressed + 1
		//		end
		//
		//		suppress_errors
		//
		//		$suppressed
		//	`,
		//	expected: 1,
		//},
		{
			name: "Rescuing an error raised by yielding to a block",
			input: `
				def say_hello(name)
					puts("Hello " + name + "!")
				end
	
				$suppressed = 0
	
				def suppress_errors
					yield
				rescue StandardError
					$suppressed = $suppressed + 1
				end

				suppress_errors { say_hello }
	
				$suppressed
			`,
			expected: 1,
		},
	}

	runVmTests(t, tests)
}
