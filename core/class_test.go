package core_test

import "testing"

func TestClass_ancestors(t *testing.T) {
	tests := []coreTestCase{
		{
			name:  "singleton Class ancestors",
			input: "Class.ancestors",
			expected: []any{
				"class:Class",
				"class:Module",
				"class:Object",
				"module:Kernel",
				"class:BasicObject",
			},
		},
		{
			name:  "Class instance ancestors",
			input: "Class.new.ancestors",
			expected: []any{
				"instance:Class",
				"class:Object",
				"module:Kernel",
				"class:BasicObject",
			},
		},
	}

	runCoreTests(t, tests)
}

func TestClass_name(t *testing.T) {
	tests := []coreTestCase{
		{
			name:     "Class",
			input:    "Class.name",
			expected: "Class",
		},
		{
			name: "Namespaced class",
			input: `
				module MyMod
					class MyClass; end
				end

				MyMod::MyClass.name
			`,
			expected: "MyMod::MyClass",
		},
	}

	runCoreTests(t, tests)
}
