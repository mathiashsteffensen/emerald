package vm

import "testing"

func TestMethodBlockArguments(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{
			name:     "block argument",
			input:    "def f(&block); block.call; end\nf { 42 }",
			expected: 42,
		},
		{
			name:     "block argument call with arguments",
			input:    "def f(&block); block.call(2) + 1; end\nf { |n| n * 3 }",
			expected: 7,
		},
		{
			name:     "block argument with a positional parameter",
			input:    "def f(value, &block); block.call(value); end\nf(5) { |n| n + 1 }",
			expected: 6,
		},
		{
			name:     "block argument with a keyword parameter",
			input:    "def f(value:, &block); block.call(value); end\nf(value: 4) { |n| n + 1 }",
			expected: 5,
		},
		{
			name:     "block argument with a default parameter",
			input:    "def f(value = 2, &block); block.call(value); end\n[f { |n| n * 4 }, f(3) { |n| n * 4 }]",
			expected: []any{8, 12},
		},
		{
			name:     "block argument definition without parentheses",
			input:    "def f value, &block; block.call(value); end\nf(3) { |n| n * 2 }",
			expected: 6,
		},
		{
			name:     "only block argument without parentheses",
			input:    "def f &block; block.call; end\nf { 11 }",
			expected: 11,
		},
		{
			name:     "yield and block argument share the block",
			input:    "def f(&block); [yield, block.call]; end\nf { 7 }",
			expected: []any{7, 7},
		},
		{
			name:     "block argument called repeatedly",
			input:    "def f(&block); block.call + block.call; end\nf { 2 }",
			expected: 4,
		},
		{
			name:     "missing block parameters become nil",
			input:    "def f(&block); block.call; end\nf { |a, b| [a, b] }",
			expected: []any{nil, nil},
		},
		{
			name:     "extra block arguments are ignored",
			input:    "def f(&block); block.call(1, 2); end\nf { |a| a }",
			expected: 1,
		},
		{
			name:     "block argument call forwards keyword arguments",
			input:    "def f(&block); block.call(value: 3); end\nf { |value:| value * 2 }",
			expected: 6,
		},
		{
			name:     "missing block binds nil",
			input:    "def f(&block); block; end\nf",
			expected: nil,
		},
		{
			name:     "empty body with a block argument returns nil",
			input:    "def f(&block); end\nf { 1 }",
			expected: nil,
		},
		{
			name:     "block argument is a Proc",
			input:    "def f(&block); block.class; end\nf { 1 }",
			expected: "class:Proc",
		},
		{
			name: "block argument captured by a nested block",
			input: `
				def apply(value, &block)
					wrap { block.call(value) }
				end
				def wrap
					yield
				end
				apply(3) { |n| n * 2 }
			`,
			expected: 6,
		},
		{
			name: "returned Proc keeps its captured context",
			input: `
				def make_proc(&block)
					block
				end
				proc = make_proc { |n| n * 4 }
				proc.call(3)
			`,
			expected: 12,
		},
		{
			name: "block argument stored in an instance variable",
			input: `
				class Greeter
					def initialize(&block)
						@block = block
					end
					def greet
						@block.call("hi")
					end
				end
				Greeter.new { |word| word + "!" }.greet
			`,
			expected: "hi!",
		},
		{
			name:     "block argument does not take a positional argument",
			input:    "def f(&block); block; end\nf(1)",
			expected: "error:ArgumentError:wrong number of arguments (given 1, expected 0)",
		},
		{
			name:     "block argument does not take a second positional argument",
			input:    "def f(value, &block); block; end\nf(1, 2) { 3 }",
			expected: "error:ArgumentError:wrong number of arguments (given 2, expected 1)",
		},
		{
			name: "block argument captures a native Enumerable block",
			input: `
				class Thing
					include Enumerable
					def each(&block)
						block.call(1)
						block.call(2)
					end
				end
				Thing.new.map { |x| x * 10 }
			`,
			expected: []any{10, 20},
		},
		{
			name:     "missing block keyword through Proc#call",
			input:    "def f(&block); block.call; end\nf { |value:| value }",
			expected: "error:ArgumentError:missing keyword: :value",
		},
		{
			name:     "missing block keyword through yield",
			input:    "def f; yield; end\nf { |value:| value }",
			expected: "error:ArgumentError:missing keyword: :value",
		},
		{
			name:     "extra block keyword is ignored",
			input:    "def f(&block); block.call(value: 1, extra: 2); end\nf { |value:| value }",
			expected: 1,
		},
		{
			name:     "block keyword argument keeps its declared slot",
			input:    "def f(&block); block.call(1, value: 2); end\nf { |a, b, value:| [a, b, value] }",
			expected: []any{1, nil, 2},
		},
		{
			name:     "block keyword argument keeps its slot after extra positional arguments",
			input:    "def f(&block); block.call(1, 2, 3, value: 4); end\nf { |a, b, value:| [a, b, value] }",
			expected: []any{1, 2, 4},
		},
		{
			name: "extra positional block arguments do not leak into locals",
			input: `
				def f(&block)
					block.call(9, 8, 7)
				end
				f { |a|
					if false
						b = 1
					end
					b
				}
			`,
			expected: nil,
		},
		{
			name:     "block keywords become a positional hash",
			input:    "def f(&block); block.call(a: 1); end\nf { |hash| hash[:a] }",
			expected: 1,
		},
		{
			name:     "block keywords append a positional hash after positional arguments",
			input:    "def f(&block); block.call(1, a: 2); end\nf { |x, hash| [x, hash[:a]] }",
			expected: []any{1, 2},
		},
		{
			name:     "extra positional hash is dropped when the block has no slot for it",
			input:    "def f(&block); block.call(1, a: 2); end\nf { |x| x }",
			expected: 1,
		},
		{
			name: "block locals after keyword parameters start as nil",
			input: `
				def f(&block)
					block.call(value: 2)
				end
				f { |value:|
					if false
						body = 1
					end
					body
				}
			`,
			expected: nil,
		},
	})
}
