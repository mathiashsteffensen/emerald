package vm

import "testing"

func TestHashLiteral(t *testing.T) {
	tests := []vmTestCase{
		{"empty hash access", "{}[:key]", nil},
		{"accessing existing symbol key", "{key: 2}[:key]", 2},
		{"accessing existing string key", `{"key" => 2}["key"]`, 2},
		{"accessing class key", "class MyClass; end; { MyClass => 2 }[MyClass]", 2},
		{"collision test", "{key: 2, \":key\" => 3}[:key]", 2}, // TODO: This should fail?
	}
	runVmTests(t, tests)
}
