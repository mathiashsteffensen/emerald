package evaluator

import "testing"

func TestEvalMethodCall(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			"Array#map#first",
			"[1, 2, 3].map { |int| int + 10 }.first",
			11,
		},
		{
			"calling a method on an identifier that evaluates to a block",
			`
			def get_arr
				[7, 9, 11]
			end

			get_arr.first`,
			7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testIntegerObject(t, testEval(tt.input), tt.expected)
		})
	}
}
