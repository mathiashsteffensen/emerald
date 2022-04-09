package evaluator

import "testing"

func TestEvalSymbolLiteral(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"normal symbol",
			":symbol",
			"symbol",
		},
		{
			"quoted symbol",
			`:"symbol-with-quotes"`,
			"symbol-with-quotes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testSymbolObject(t, testEval(tt.input), tt.expected)
		})
	}
}
