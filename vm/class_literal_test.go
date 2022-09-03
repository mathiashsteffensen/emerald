package vm

import "testing"

func TestClassLiteral(t *testing.T) {
	tests := []vmTestCase{
		{
			name:     "empty class",
			input:    "class MyClass; end",
			expected: nil,
		},
		{
			name: "class with method",
			input: `class MyClass
				def my_method; end
			end`,
			expected: ":my_method",
		},
		{
			name: "class with arbitrary last expression",
			input: `class MyClass
				def my_method; end
				def my_other_method; end

				69
			end`,
			expected: 69,
		},
		{
			name: "class with included module",
			input: `
			module MyMod; end

			class MyClass
				include(MyMod)
			end`,
			expected: "class:MyClass",
		},
		{
			name: "namespaced class",
			input: `
			module MyMod
				class MyClass; end
			end
		
			MyMod::MyClass`,
			expected: "class:MyClass",
		},
		{
			name: "inheriting class",
			input: `
				class Parent
					def number
						5
					end
				end
				
				class Child < Parent; end

				Child.new.number
			`,
			expected: 5,
		},
	}

	runVmTests(t, tests)
}
