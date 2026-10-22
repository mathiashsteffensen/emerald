package vm

import "testing"

func TestEmptyBoundRescueReturnsNil(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "typed rescue binding", input: `
			def f
				raise "boom"
			rescue StandardError => error
			end
			f
		`, expected: nil},
		{name: "default rescue binding", input: "def f; raise \"boom\"\nrescue => error\nend; f", expected: nil},
		{name: "existing local binding", input: `
			def f
				error = nil
				raise "boom"
			rescue StandardError => error
			end
			f
		`, expected: nil},
		{name: "empty yielded rescue binding", input: `
			[1].map do |n|
				raise "boom"
			rescue StandardError => error
			end
		`, expected: []any{nil}},
	})
}

func TestRescueVariableBinding(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{input: `def f; raise "boom"; rescue RuntimeError, TypeError => error; error.to_s; end; f`, expected: "boom"},
		{input: `def f; error = nil; raise "boom"; rescue StandardError => error; error.to_s; end; f`, expected: "boom"},
		{input: "def f; raise \"boom\"\nrescue => error\nerror.to_s; end; f", expected: "boom"},
	})
}

func TestErrorHandling(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "calling method with wrong amount of arguments",
			input: `
			def say_hello(name)
				puts("Hello " + name + "!")
			end
		
			say_hello
			`,
			expected: "error:ArgumentError:wrong number of arguments (given 0, expected 1)",
		},
		{
			name: "Rescuing an error",
			input: `
				def say_hello(name)
					puts("Hello " + name + "!")
				end
		
				$suppressed = 0
		
				def suppress_argument_errors
					say_hello
				rescue ArgumentError
					$suppressed = $suppressed + 1
				end
		
				suppress_argument_errors
		
				$suppressed
			`,
			expected: 1,
		},
		{
			name: "Rescuing an error by specifying subclass",
			input: `
				def say_hello(name)
					puts("Hello " + name + "!")
				end
		
				$suppressed = 0
		
				def suppress_errors
					say_hello
				rescue StandardError
					$suppressed = $suppressed + 1
				end
		
				suppress_errors
		
				$suppressed
			`,
			expected: 1,
		},
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
		{
			name: "Rescue clause in block",
			input: `
				def say_hello(name)
					puts("Hello " + name + "!")
				end

				def yielder
					yield
				end
	
				$suppressed = 0

				yielder do 
					say_hello
				rescue StandardError
					$suppressed = $suppressed + 1
				end
	
				$suppressed
			`,
			expected: 1,
		},
	}

	runVmTests(t, tests)
}

func TestRescueThroughNativeIterator(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "through each", input: `
			def f
				result = []
				[1, 2].each { |x| result.push(x); raise "a" }
				result.push("after")
			rescue
				result.push("rescued")
			end
			f
		`, expected: []any{1, "rescued"}},
		{name: "through map", input: `
			def f
				result = []
				[1, 2].map { |x| result.push(x); raise "a" }
				result.push("after")
			rescue
				result.push("rescued")
			end
			f
		`, expected: []any{1, "rescued"}},
		{name: "through a bytecode method that yields", input: `
			def twice
				yield 1
				yield 2
			end

			def f
				result = []
				twice { |x| result.push(x); raise "a" }
				result.push("after")
			rescue
				result.push("rescued")
			end
			f
		`, expected: []any{1, "rescued"}},
		{name: "through nested iterators", input: `
			def f
				result = []
				[1, 2].each { |x| [3, 4].map { |y| result.push(y); raise "a" } }
			rescue
				result.push("rescued")
			end
			f
		`, expected: []any{3, "rescued"}},
		{name: "continues after the call", input: `
			def f
				[1, 2].each { |x| raise "a" }
			rescue
				"rescued"
			end
			[f, f, 1 + 1]
		`, expected: []any{"rescued", "rescued", 2}},
		{name: "restores self", input: `
			class Foo
				def f
					[1].each { raise "a" }
				rescue
					name
				end

				def name
					"foo"
				end
			end
			Foo.new.f
		`, expected: "foo"},
		{name: "clears $!", input: `
			def f
				[1].each { raise "a" }
			rescue
			end
			f
			$!
		`, expected: nil},
	})
}
