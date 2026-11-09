package vm

import "testing"

// The lexer has no token for these operator characters, so it hands them to
// the parser as ILLEGAL; they must still work as method names.
func TestOperatorCharactersWithoutTokenAsMethodNames(t *testing.T) {
	setup := `
		class Foo
			def %(o)
				o * 10
			end

			def ~
				7
			end

			def ^(o)
				o + 1
			end
		end

		f = Foo.new
	`

	tests := []vmTestCase{
		{name: "def % called with f.%(4)", input: "f.%(4)", expected: 40},
		{name: "def ~ called with f.~", input: "f.~", expected: 7},
		{name: "def ^ called with f.^(1)", input: "f.^(1)", expected: 2},
	}

	runVmTests(t, tests, setup)
}
