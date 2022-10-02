package vm

import "testing"

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"", `"monkey"`, "monkey"},
		{"", `"mon" + "key"`, "monkey"},
		{"", `"mon" + "key" + "banana"`, "monkeybanana"},
		{"", `placeholder = "template"; "This is a #{placeholder}"`, "This is a template"},
	}
	runVmTests(t, tests)
}
