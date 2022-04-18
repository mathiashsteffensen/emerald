package vm

import "testing"

func TestMethodCall(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "built in instance method",
			input:    `"string".to_sym`,
			expected: ":string",
		},
		{
			name: "creating instance and instance method",
			input: `
			class Greeter
				def hello
					"hello"
				end
			end

			instance = Greeter.new

			instance.hello
			`,
			expected: "hello",
		},
		{
			name: "creating class with static method",
			input: `
			class Greeter
				class << self
					def hello
						"hello"
					end
				end
			end

			Greeter.hello
			`,
			expected: "hello",
		},
		{
			name:     "passing a block to a builtin method",
			input:    "[0,1,2].map { |i| i + 2 }",
			expected: []any{2, 3, 4},
		},
		{
			name: "calling a top level method within a block passed to a builtin method",
			input: `
			def add_two(n); n + 2; end
			[0,1,2].map { |i| add_two(i) }`,
			expected: []any{2, 3, 4},
		},
		{
			name: "calling a top level method within a block passed to a builtin method",
			input: `
			def add_two(n); n + 2; end
			[0,1,2].map { |i| add_two(i) }.sum`,
			expected: 9,
		},
		{
			name: "calling a method with a receiver within a block passed to a builtin method",
			input: `
			class Math
				def add_two(n); n + 2; end
				class << self
					def instance; new; end
				end
			end
			[0,1,2].map { |i| Math.instance.add_two(i) }.sum`,
			expected: 9,
		},
		{
			name: "closure test",
			input: `
				class MyClass
					[:one, :two, :three].each do |lvl|
						define_method(lvl) do |other_val|
							lvl
						end
					end
				end

				MyClass.new.one(:two)
			`,
			expected: ":one",
		},
		{
			input:    "Object.new",
			expected: "instance:Object",
		},
	}

	runVmTests(t, tests)
}
