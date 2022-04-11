package vm

import "testing"

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"", `"monkey"`, "monkey"},
		{"", `"mon" + "key"`, "monkey"},
		{"", `"mon" + "key" + "banana"`, "monkeybanana"},
	}
	runVmTests(t, tests)
}
