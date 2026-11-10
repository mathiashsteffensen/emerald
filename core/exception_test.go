package core_test

import (
	"emerald/core"
	"fmt"
	"strings"
	"testing"
)

func TestExceptionClassNameWithoutInstance(t *testing.T) {
	err := &core.ExceptionInstance{}
	if got := err.ClassName(); got != "Exception" {
		t.Fatalf("class name = %q, want Exception", got)
	}
}

func TestException_kind_of(t *testing.T) {
	tests := []coreTestCase{}

	rt := core.NewRuntime()
	rt.Init()

	for className := range rt.Object.NamespaceDefinitions() {
		if strings.Contains(className, "Error") {
			tests = append(tests, coreTestCase{
				name:     fmt.Sprintf("%s#kind_of?(Exception) == true", className),
				input:    fmt.Sprintf("%s.new.kind_of?(Exception)", className),
				expected: true,
			})
		}
	}

	runCoreTests(t, tests)
}

func TestException_initialize(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "subclass with message",
			input:    `MyError.new("boom").message`,
			expected: "boom",
		},
		{
			name:     "subclass without message",
			input:    `MyError.new.message == ""`,
			expected: true,
		},
		{
			name:     "too many arguments",
			input:    `MyError.new("a", "b")`,
			expected: "error:ArgumentError:wrong number of arguments (given 2, expected 0..1)",
		},
		{
			name:     "keywords are ignored",
			input:    `MyError.new(foo: 1).class == MyError`,
			expected: true,
		},
		{
			name:     "built-in instance keeps its message",
			input:    `error = ArgumentError.new("x"); error.initialize("y"); error.message`,
			expected: "x",
		},
		{
			name:     "NameError subclass accepts a name",
			input:    `MyNameError.new("m", :foo).message`,
			expected: "m",
		},
		{
			name:     "NameError subclass with too many arguments",
			input:    `MyNameError.new("m", :foo, 1)`,
			expected: "error:ArgumentError:wrong number of arguments (given 3, expected 0..2)",
		},
		{
			name:     "NoMethodError subclass accepts a name and arguments",
			input:    `MyNoMethodError.new("m", :foo, [1]).message`,
			expected: "m",
		},
		{
			name:     "NoMethodError subclass with too many arguments",
			input:    `MyNoMethodError.new("m", :foo, [1], true, 1)`,
			expected: "error:ArgumentError:wrong number of arguments (given 5, expected 0..4)",
		},
	}

	prelude := `
		class MyError < StandardError; end
		class MyNameError < NameError; end
		class MyNoMethodError < NoMethodError; end
	`

	runCoreTests(t, tests, prelude)
}
