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
				module MyMod
					class << self
						define_method(:hello) { "Hello" }
					end
				end
		
				MyMod.hello
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

func TestModule_attr_accessor(t *testing.T) {
	tests := []coreTestCase{
		{
			name: "accessor",
			input: `
				class MyClass
					attr_accessor :hello
				end

				c = MyClass.new

				c.hello = "Hello"
				c.hello
			`,
			expected: "Hello",
		},
	}

	runCoreTests(t, tests)
}
