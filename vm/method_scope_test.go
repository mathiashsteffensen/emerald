package vm

import "testing"

func TestMethodScopeGate(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "method local does not overwrite top-level local",
			input: `
				left = 1
				def m
					left = 9
					left
				end
				[m, left]
			`,
			expected: []any{9, 1},
		},
		{
			name: "method body does not read top-level local",
			input: `
				def x; "method"; end
				x = 1
				def m; x; end
				m
			`,
			expected: "method",
		},
		{
			name: "def inside class",
			input: `
				x = 5
				class Foo
					def bar
						x = 7
						x
					end
				end
				[Foo.new.bar, x]
			`,
			expected: []any{7, 5},
		},
		{
			name: "def inside module",
			input: `
				x = 5
				module Bar
					def baz
						x = 7
					end
				end
				class Baz
					include Bar
				end
				Baz.new.baz
				x
			`,
			expected: 5,
		},
		{
			name: "nested def does not see outer method local",
			input: `
				def outer
					z = 1
					def inner
						z = 2
					end
					inner
					z
				end
				outer
			`,
			expected: 1,
		},
		{
			name: "block in method does not see top-level local",
			input: `
				w = 1
				def m
					[1].each { w = 2 }
				end
				m
				w
			`,
			expected: 1,
		},
		{
			name:     "top-level block captures top-level local",
			input:    "y = 2; [1, 2].each { |i| y += i }; y",
			expected: 5,
		},
		{
			name: "block in method captures method local",
			input: `
				class Foo
					def bar
						x = 7
						[1].each { |i| x += i }
						x
					end
				end
				Foo.new.bar
			`,
			expected: 8,
		},
		{
			name: "globals resolve through method scope",
			input: `
				$g = 3
				def gm
					$g = $g + 1
				end
				gm
				$g
			`,
			expected: 4,
		},
	}

	runVmTests(t, tests)
}
