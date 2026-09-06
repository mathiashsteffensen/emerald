package vm

import "testing"

func TestCapturedLocalMutation(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "captured sum", input: `
			def sum(values)
				total = 0
				values.each { |n| total += n }
				total
			end
			[sum([1, 2, 3]), sum([4, 5])]
		`, expected: []any{6, 9}},
		{name: "nested capture", input: `
			def f
				n = 0
				[1, 2].each { |a| [3, 4].each { |b| n += a + b } }
				n
			end
			f
		`, expected: 20},
		{name: "escaped sibling closures share updates", input: `
			class Counter
				[10].each do |n|
					define_method(:increment) { n += 1 }
					define_method(:value) { n }
					n = 20
				end
			end
			counter = Counter.new
			[counter.increment, counter.increment, counter.value]
		`, expected: []any{21, 22, 22}},
	})
}

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
			input:    "[0,1,2].map { |i| i }",
			expected: []any{0, 1, 2},
		},
		{
			name:     "passing a block to a builtin method v2",
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
		{
			input: `
			module MyMod
				def hello
					"Hello"
				end
			end

			class MyClass
				include(MyMod)
			end

			MyClass.new.hello`,
			expected: "Hello",
		},
		{
			name: "keyword arguments",
			input: `
					def m(n:, o:)
						n
					end
					m(o: 3, n: 1)
				`,
			expected: 1,
		},
	}

	runVmTests(t, tests)
}
