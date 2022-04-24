package core_test

import "testing"

func TestModule_define_method(t *testing.T) {
	tests := []coreTestCase{
		{
			name: "defining an instance method",
			input: `
				class MyClass
					define_method(:hello) { "Hello" }
				end

				MyClass.new.hello
			`,
			expected: "Hello",
		},
		{
			name: "defining a static method",
			input: `
				class MyClass
					class << self
						define_method(:hello) { "Hello" }
					end
				end

				MyClass.hello
			`,
			expected: "Hello",
		},
		{
			name: "defining a method inside a block in a class",
			input: `
				METHODS = [:hello]

				class MyClass
					10.times do
						METHODS.each do |method|
							define_method(method) { "Hello" }
						end
					end
				end

				MyClass.new.hello
			`,
			expected: "Hello",
		},
		{
			name: "defining a method inside a block in a module",
			input: `
				METHODS = [:hello]

				module MyMod
					10.times do
						METHODS.each do |method|
							define_method(method) { "Hello" }
						end
					end
				end

				class MyClass
					include(MyMod)
				end

				MyClass.new.hello
			`,
			expected: "Hello",
		},
	}

	runCoreTests(t, tests)
}
