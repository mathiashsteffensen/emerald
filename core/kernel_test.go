package core_test

import "testing"

func TestKernel_class(t *testing.T) {
	tests := []coreTestCase{
		{
			input:    "Object.new.class",
			expected: "class:Object",
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
		{
			name:     "with explicit module receiver",
			input:    `Kernel.puts(nil)`,
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
	}

	runCoreTests(t, tests)
}
