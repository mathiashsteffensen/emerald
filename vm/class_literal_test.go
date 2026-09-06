package vm

import "testing"

func TestInheritedClassMethods(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "inherit and override", input: `
			class Parent
				class << self
					def value; 7; end
				end
			end
			class Child < Parent; end
			class Other < Parent
				class << self
					def value; 9; end
				end
			end
			[Parent.value, Child.value, Other.value, Child.class]
		`, expected: []any{7, 7, 9, "class:Class"}},
		{name: "parent method added later", input: `
			class Parent; end
			class Child < Parent; end
			Child.new
			class Parent
				class << self
					def value; 7; end
				end
			end
			Child.value
		`, expected: 7},
	})
}

func TestSingletonClassReceiver(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "object receiver and outer self", input: `
			obj = Object.new
			outer = self
			class << obj
				def value; 7; end
			end
			[obj.value, self == outer]
		`, expected: []any{7, true}},
		{name: "empty singleton body", input: "class << self; end", expected: nil},
	})
}

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
				module MyOtherMod
					class MyClass; end
				end
			end
		
			MyMod::MyOtherMod::MyClass`,
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
