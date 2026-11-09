package vm

import "testing"

func TestMethodDefaultArguments(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{
			name:     "omitted trailing default",
			input:    "def f(a, b = 10); a + b; end\nf(3)",
			expected: 13,
		},
		{
			name:     "supplied trailing default",
			input:    "def f(a, b = 10); a + b; end\nf(3, 4)",
			expected: 7,
		},
		{
			name:     "default uses an earlier parameter",
			input:    "def f(a, b = a + 1); b; end\nf(3)",
			expected: 4,
		},
		{
			name:     "several defaults",
			input:    "def f(a = 1, b = 2, c = 3); [a, b, c]; end\nf(9)",
			expected: []any{9, 2, 3},
		},
		{
			name:     "no arguments uses every default",
			input:    "def f(a = 1, b = 2, c = 3); [a, b, c]; end\nf",
			expected: []any{1, 2, 3},
		},
		{
			name:     "required parameter after defaults takes the last argument",
			input:    "def f(a = 1, b = 2, c); [a, b, c]; end\nf(10, 20)",
			expected: []any{10, 2, 20},
		},
		{
			name:     "required parameter after default takes the only argument",
			input:    "def f(a = 1, b); [a, b]; end\nf(5)",
			expected: []any{1, 5},
		},
		{
			name:     "string default",
			input:    "def f(greeting = \"Hello\"); greeting; end\nf",
			expected: "Hello",
		},
		{
			name:     "nil default",
			input:    "def f(value = nil); value; end\nf",
			expected: nil,
		},
		{
			name:     "explicit nil overrides default",
			input:    "def f(value = 5); value; end\nf(nil)",
			expected: nil,
		},
		{
			name:     "explicit false overrides default",
			input:    "def f(value = true); value; end\nf(false)",
			expected: false,
		},
		{
			name: "unused default is not evaluated",
			input: `
				$effects = 0
				def bump
					$effects += 1
					9
				end
				def f(value = bump)
					value
				end
				[f(1), $effects]
			`,
			expected: []any{1, 0},
		},
		{
			name: "used default is evaluated once",
			input: `
				$effects = 0
				def bump
					$effects += 1
					9
				end
				def f(value = bump)
					value
				end
				[f, $effects]
			`,
			expected: []any{9, 1},
		},
		{
			name: "default calls a method on self",
			input: `
				class Widget
					def default_size
						3
					end
					def size(value = default_size)
						value
					end
				end
				Widget.new.size
			`,
			expected: 3,
		},
		{
			name:     "default calls a top level method",
			input:    "def seven; 7; end\ndef f(value = seven); value; end\nf",
			expected: 7,
		},
		{
			name:     "default contains a block",
			input:    "def f(values = [1, 2].map { |n| n * 2 }); values; end\nf",
			expected: []any{2, 4},
		},
		{
			name:     "default with keyword argument only call",
			input:    "def f(value = 1, suffix:); [value, suffix]; end\nf(suffix: 2)",
			expected: []any{1, 2},
		},
		{
			name:     "default with positional and keyword arguments",
			input:    "def f(value = 1, suffix:); [value, suffix]; end\nf(5, suffix: 2)",
			expected: []any{5, 2},
		},
		{
			name: "default in constructor",
			input: `
				class Box
					def initialize(width = 2, height = 3)
						@width = width
						@height = height
					end
					def area
						@width * @height
					end
				end
				[Box.new.area, Box.new(4).area, Box.new(4, 5).area]
			`,
			expected: []any{6, 12, 20},
		},
		{
			name: "default in singleton method",
			input: `
				class Widget
					class << self
						def build(size = 7)
							size
						end
					end
				end
				[Widget.build, Widget.build(2)]
			`,
			expected: []any{7, 2},
		},
		{
			name: "default through operator dispatch",
			input: `
				class Box
					def +(other, extra = 10)
						other + extra
					end
				end
				Box.new + 2
			`,
			expected: 12,
		},
		{
			name:     "default in method definition without parentheses",
			input:    "def f a, b = 2; a + b; end\n[f(1), f(1, 3)]",
			expected: []any{3, 4},
		},
		{
			name: "default error is not caught by the method rescue",
			input: `
				def explode(value = raise("boom"))
					value
				rescue StandardError
					"method rescued"
				end
				def caller
					explode
				rescue StandardError
					"caller rescued"
				end
				caller
			`,
			expected: "caller rescued",
		},
		{
			name:     "too few arguments",
			input:    "def f(a, b = 2); a; end\nf",
			expected: "error:ArgumentError:wrong number of arguments (given 0, expected 1..2)",
		},
		{
			name:     "too many arguments",
			input:    "def f(a, b = 2); a; end\nf(1, 2, 3)",
			expected: "error:ArgumentError:wrong number of arguments (given 3, expected 1..2)",
		},
		{
			name:     "all parameters defaulted arity",
			input:    "def f(a = 1, b = 2); a; end\nf(1, 2, 3)",
			expected: "error:ArgumentError:wrong number of arguments (given 3, expected 0..2)",
		},
		{
			name:     "method without defaults keeps its exact arity message",
			input:    "def f(a); a; end\nf(1, 2)",
			expected: "error:ArgumentError:wrong number of arguments (given 2, expected 1)",
		},
	})
}
