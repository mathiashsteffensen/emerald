package vm

import "testing"

func TestIfModifier(t *testing.T) {
	setup := `
		$i = 1
		
		def call
			$i = $i + 1
		end
	`

	tests := []vmTestCase{
		{
			input: `
				call if true && false

				$i
			`,
			expected: 1,
		},
		{
			input: `
				call if 1 + 1 == 2

				$i
			`,
			expected: 2,
		},
		{
			input: `
				def call?(should_we)
					2 * 5 == 10 if should_we
				end

				call unless call? false
				call if call? 2 * 9 + 17 == 35
				call if call? false
				call if call? true ||
				
									false
				$i
			`,
			expected: 4,
		},
	}

	runVmTests(t, tests, setup)
}

func TestModifierAfterCallArguments(t *testing.T) {
	runVmTests(t, []vmTestCase{
		{name: "if false", input: `raise "x" if false; :reached`, expected: ":reached"},
		{name: "if true", input: `raise "x" if true`, expected: "error:RuntimeError:x"},
		{name: "unless true", input: `raise "x" unless true; :reached`, expected: ":reached"},
		{name: "after keyword argument", input: `def log(message, level:); [message, level]; end; log "x", level: 1 if true`, expected: []any{"x", 1}},
	})
}
