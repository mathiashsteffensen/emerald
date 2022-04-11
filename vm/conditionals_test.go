package vm

import (
	"testing"
)

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"", "if true; 10; end", 10},
		{"", "if true; 10; else; 20; end", 10},
		{"", "if false; 10; else; 20; end", 20},
		{"", "if 1; 10; end", 10},
		{"", "if 1 < 2; 10; end", 10},
		{"", "if 1 < 2; 10; else; 20; end", 10},
		{"", "if 1 > 2; 10; else; 20; end", 20},
		{"", "if 1 > 2; 10; end", nil},
		{"", "if (if false 10 end) 10 else 20 end", 20},
	}
	runVmTests(t, tests)
}
