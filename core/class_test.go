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
			name: "Namespaced class in module",
			input: `
				module MyMod
					class MyClass; end
				end

				MyMod::MyClass.name
			`,
			expected: "MyMod::MyClass",
		},
		{
			name: "Namespaced class in class",
			input: `
				class MyClass
					class MyOtherClass; end
				end

				MyClass::MyOtherClass.name
			`,
			expected: "MyClass::MyOtherClass",
		},
	}

	runCoreTests(t, tests)
}

func TestClass_new_instance(t *testing.T) {
	before := `
		def create_class(kwargs:)
			module MyMod
				class MyClass
					if kwargs
						def initialize(value:)
							@value = value
						end
					else
						def initialize(value)
							@value = value
						end
					end

					def value
						@value
					end
				end
			end
		end
	`

	tests := []coreTestCase{
		{
			name: "positional argument provided",
			input: `
								create_class(kwarfs: false)

								instance = MyMod::MyClass.new(2)
								instance.value
							`,
			expected: 2,
		},
		{
			name: "positional argument omitted",
			input: `
								create_class(kwargs: false)

								instance = MyMod::MyClass.new()
								instance.value
							`,
			expected: "error:ArgumentError:wrong number of arguments (given 0, expected 1)",
		},
		{
			name: "keyword argument provided",
			input: `
						create_class(kwargs: true)

						instance = MyMod::MyClass.new(value: 3)
						instance.value
					`,
			expected: 3,
		},
		{
			name: "keyword argument omitted",
			input: `
				create_class(kwargs: true)

				instance = MyMod::MyClass.new
				instance.value
			`,
			expected: "error:ArgumentError:missing keyword: :value",
		},
	}

	runCoreTests(t, tests, before)
}
