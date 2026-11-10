package core_test

import "testing"

func TestProc(t *testing.T) {
	runCoreTests(t, []coreTestCase{
		{
			name:     "new wraps the call's block",
			input:    "Proc.new { 42 }.call",
			expected: 42,
		},
		{
			name:     "new returns a Proc",
			input:    "Proc.new { 1 }.class",
			expected: "class:Proc",
		},
		{
			name:     "new without a block",
			input:    "Proc.new",
			expected: "error:ArgumentError:tried to create Proc object without a block",
		},
		{
			name:     "new with arguments rejects them",
			input:    "Proc.new(1) { 2 }",
			expected: "error:ArgumentError:wrong number of arguments (given 1, expected 0)",
		},
		{
			name:     "subclass new returns an instance of the subclass",
			input:    "class MyProc < Proc; end; MyProc.new { 1 }.class",
			expected: "class:MyProc",
		},
		{
			name: "subclass initialize receives the block and arguments",
			input: `
				class MyProc < Proc
					def initialize(value, &block)
						@value = value
						@block_class = block.class
					end
					def value
						@value
					end
					def block_class
						@block_class
					end
				end
				proc = MyProc.new(7) { 9 }
				[proc.class == MyProc, proc.value, proc.block_class == Proc, proc.call]
			`,
			expected: []any{true, 7, true, 9},
		},
		{
			name:     "subclass without initialize inherits Proc's arity check",
			input:    "class MyProc < Proc; end; MyProc.new(1) { 2 }",
			expected: "error:ArgumentError:wrong number of arguments (given 1, expected 0)",
		},
		{
			name:     "call a captured block",
			input:    "def f(&block); block.call(2); end; f { |n| n * 3 }",
			expected: 6,
		},
		{
			name:     "captured block class",
			input:    "def f(&block); block.class == Proc; end; f { 1 }",
			expected: true,
		},
		{
			name:     "call without a captured block",
			input:    "Proc.new.call",
			expected: "error:ArgumentError:tried to create Proc object without a block",
		},
	})
}
