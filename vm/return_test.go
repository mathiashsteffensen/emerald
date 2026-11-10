package vm

import (
	"emerald/compiler"
	"emerald/core"
	"testing"
)

func TestReturnFromBlock(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "in a for loop body in a method",
			input: `
				def find
					for x in [1, 2, 3]
						return x * 100 if x == 2
						$seen << x
					end
					:not_here
				end
				$seen = []
				[find, $seen]
			`,
			expected: []any{200, []any{1}},
		},
		{
			name: "in a for loop nested in a while loop in a for loop",
			input: `
				def find
					for a in [1, 2]
						i = 0
						while i < 3
							for b in [10, 20]
								return [a, i, b] if a == 2 && i == 1 && b == 20
							end
							i += 1
						end
					end
					:not_here
				end
				find
			`,
			expected: []any{2, 1, 20},
		},
		{
			name:     "through native each",
			input:    "def find; [1, 2, 3].each { |x| return x if x == 2 }; :not_here; end; find",
			expected: 2,
		},
		{
			name:     "through native map and a nested block",
			input:    "def find; [1, 2].map { |x| [3].each { return x * 10 } }; :not_here; end; find",
			expected: 10,
		},
		{
			name:     "through a method that yields",
			input:    "def twice; yield; yield; :not_here; end; def find; twice { return 7 }; :not_here; end; find",
			expected: 7,
		},
		{
			name:     "without a value",
			input:    "def find; [1].each { return }; :not_here; end; find",
			expected: nil,
		},
		{
			name:     "past a rescue",
			input:    "def find; [1].each { return 5 }; :not_here; rescue => e; :rescued; end; find",
			expected: 5,
		},
		{
			name:     "from a block argument called while its method runs",
			input:    "def run(&b); b.call; :not_here; end; def find; run { return 4 }; :not_here; end; find",
			expected: 4,
		},
		{
			name: "from the right invocation of a recursive method",
			input: `
				def walk(n)
					[1].each do
						return n if n == 0
						inner = walk(n - 1)
						return [n, inner]
					end
					:not_here
				end
				walk(2)
			`,
			expected: []any{2, []any{1, 0}},
		},
		{
			name:     "in a block of a method defined in a block",
			input:    "[1].each { def find; [1].each { return :found }; :not_here; end }; find",
			expected: ":found",
		},
		{
			name:     "after the method returned",
			input:    "def keep(&b); b; end; def make; keep { return 1 }; end; make.call",
			expected: "error:LocalJumpError:unexpected return",
		},
		{
			name:     "after the method returned is rescuable",
			input:    "def keep(&b); b; end; def make; keep { return 1 }; end; def run; make.call; rescue LocalJumpError => e; e.message; end; run",
			expected: "unexpected return",
		},
		{
			name:     "break still leaves the call",
			input:    "def find; [1, 2, 3].each { |x| break x * 2 if x == 2 }; end; find",
			expected: 4,
		},
		{
			name:     "next still leaves the block",
			input:    "def find; [1, 2].map { |x| next x * 3 }; end; find",
			expected: []any{3, 6},
		},
	}

	runVmTests(t, tests)
}

func TestReturnFromBlockAtTheTopLevelEndsTheScript(t *testing.T) {
	rt := core.NewRuntime()
	rt.Init()

	bc := compiler.Compile("test", "$x = 1; [1].each { return }; $x = 2", rt)
	New("test", bc, rt).Run()

	ensureNoExceptionUnlessExpected(t, nil, rt)

	if err := testIntegerObject(1, rt.Heap.GetGlobalVariableString("$x")); err != nil {
		t.Errorf("the script went on after the return: %s", err)
	}
}

func TestMethodLocalsStartAsNil(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "after a for loop with no iterations",
			input: `
				def noise(a, b, c); [a, b, c, "s", "t", "u"].map { |v| v }; end
				def zero; for x in []; y = 1; end; [x, y]; end
				noise(:leak1, :leak2, :leak3)
				zero
			`,
			expected: []any{nil, nil},
		},
		{
			name: "assigned in a branch not taken",
			input: `
				def noise(a, b, c); [a, b, c].map { |v| v }; end
				def m; if false; y = 1; end; y; end
				noise(:leak1, :leak2, :leak3)
				m
			`,
			expected: nil,
		},
	}

	runVmTests(t, tests)
}

func TestReturnFromDefineMethodBlock(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "called directly",
			input:    "class Foo; define_method(:bar) { |x| return x * 2; :not_here }; end; Foo.new.bar(4)",
			expected: 8,
		},
		{
			name: "called as an operator",
			input: `
				class Foo
					def ==(other)
						[1, 2].each { |x| return :equal if x == other }
						:not_equal
					end
				end
				[Foo.new == 2, Foo.new == 3]
			`,
			expected: []any{":equal", ":not_equal"},
		},
		{
			name:     "called as to_s from interpolation",
			input:    `class Foo; define_method(:to_s) { return "foo"; "not_here" }; end; "#{Foo.new}"`,
			expected: "foo",
		},
		{
			name: "called as initialize from new",
			input: `
				class Foo
					define_method(:initialize) { |x| @x = x; return; @x = :not_here }
					define_method(:x) { @x }
				end
				$after = :after
				[Foo.new(5).x, $after]
			`,
			expected: []any{5, ":after"},
		},
	}

	runVmTests(t, tests)
}

func TestBlockLocalsStartAsNil(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "with fewer arguments than parameters",
			input:    "def noise; [:leak1, :leak2, :leak3].map { |v| v }; end; def call; yield 1; end; noise; call { |a, b| y = 1 if false; [a, b, y] }",
			expected: []any{1, nil, nil},
		},
		{
			name:     "with no arguments",
			input:    "def noise; [:leak1, :leak2, :leak3].map { |v| v }; end; def call; yield; end; noise; call { |a| y = 1 if false; [a, y] }",
			expected: []any{nil, nil},
		},
	}

	runVmTests(t, tests)
}
