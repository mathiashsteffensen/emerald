package core_test

import "testing"

func TestKernel_require_relative(t *testing.T) {
	tests := []coreTestCase{
		{
			input: `
				require_relative("fixtures/require_test")
				MyClass.run_require
				require_relative("../spec/fixtures/require_test")
			`,
			expected: true,
		},
		{
			name: "requiring same file twice",
			input: `
				require_relative("fixtures/require_test")
				require_relative("fixtures/require_test")
			`,
			expected: false,
		},
		// TODO: Uncomment when exception handling has been implemented
		//{
		//	name:     "when file doesn't exist",
		//	input:    `require_relative("../lib/main")`,
		//	expected: "error:LoadError:cannot load such file -- /home/mathias/code/emerald/lib/main",
		//},
	}

	runCoreTests(t, tests)
}

func TestKernel_class(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "Object.new.class",
			expected: "class:Object",
		},
		{
			input:    "Object.class",
			expected: "class:Class",
		},
	}

	runCoreTests(t, tests)
}

func TestKernel_kind_of(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "when self is instance of class",
			input:    "Object.new.kind_of?(Object)",
			expected: true,
		},
		{
			name:     "when self is instance of subclass",
			input:    "String.new.kind_of?(Object)",
			expected: true,
		},
		{
			name:     "singleton is instance of Class",
			input:    "String.kind_of?(Class)",
			expected: true,
		},
		{
			name: "when passed a not included module",
			input: `
			module MyMod; end
			class MyClass; end
			MyClass.new.kind_of?(MyMod)`,
			expected: false,
		},
		{
			name: "when passed an included module",
			input: `
			module MyMod; end
			class MyClass
				include(MyMod)
			end
			MyClass.new.kind_of?(MyMod)`,
			expected: true,
		},
		{
			name:     "when passed wrong type of arg",
			input:    "String.kind_of?(23)",
			expected: "error:TypeError:class or module required",
		},
	}

	runCoreTests(t, tests)
}

func TestKernel_puts(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "with a single string argument",
			input:    `puts("Hello World!")`,
			expected: nil,
		},
		{
			name:     "with multiple string arguments",
			input:    `puts("Hello", "World!")`,
			expected: nil,
		},
		{
			name:     "with nil argument",
			input:    `puts(nil)`,
			expected: nil,
		},
	}

	runCoreTests(t, tests)
}

func TestKernel_include(t *testing.T) {
	tests := []coreTestCase{
		{
			name: "in main object",
			input: `module MyMod
				def hello; "Hello"; end
			end
			include(MyMod)
			hello`,
			expected: "Hello",
		},
		{
			name: "in custom class",
			input: `module MyMod
				def hello; "Hello"; end
			end
		
			class MyClass
				include(MyMod)
			end
		
			MyClass.new.hello`,
			expected: "Hello",
		},
		{
			name: "in custom static class",
			input: `module MyMod
				def hello; "Hello"; end
			end
			
			class MyClass
				class << self
					include(MyMod)
				end
			end

			MyClass.hello`,
			expected: "Hello",
		},
	}

	runCoreTests(t, tests)
}
