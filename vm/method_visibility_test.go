package vm

import "testing"

func TestMethodVisibility(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `puts "Hello"`,
			expected: nil,
		},
		{
			input:    `Kernel.puts "Hello"`,
			expected: nil,
		},
		{
			input:    `Object.puts "Hello"`,
			expected: "error:NoMethodError:private method `puts' called for Object:Class",
		},
		{
			input: `
				class C
					private def call; 1; end
				end

				C.new.call
			`,
			expected: "error:NoMethodError:private method `call' called for C",
		},
		{
			input: `
				class C
					private

					def call; 1; end
				end

				C.new.call
			`,
			expected: "error:NoMethodError:private method `call' called for C",
		},
	}

	runVmTests(t, tests)
}
