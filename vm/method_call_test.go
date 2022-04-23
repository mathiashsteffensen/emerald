package vm

import "testing"

func TestMethodCall(t *testing.T) {
	tests := []vmTestCase{
		{
			name: "closure test",
			input: `
				class MyClass
					[:one, :two, :three].each do |lvl|
						define_method(lvl) do |other_val|
							lvl
						end
					end
				end

				MyClass.new.one(:two)
			`,
			expected: ":one",
		},
		{
			input:    "Object.new",
			expected: "instance:Object",
		},
		{
			input: `
			module MyMod
				def hello
					"Hello"
				end
			end

			class MyClass
				include(MyMod)
			end

			MyClass.new.hello`,
			expected: "Hello",
		},
	}

	runVmTests(t, tests)
}
