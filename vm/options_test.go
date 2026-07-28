package vm

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/object"
	"testing"
)

func TestNewUsesEmptyInputs(t *testing.T) {
	rt := core.NewRuntime()
	rt.Init()

	New("test", &bytecode.Bytecode{}, rt)

	assertStringArray(t, rt.Heap.GetGlobalVariableString("$LOAD_PATH"))
	assertStringArray(t, rt.Heap.GetGlobalVariableString("$*"))

	argv, err := getConst(rt.MainObject, "ARGV", rt)
	if err != nil {
		t.Fatal(err)
	}
	assertStringArray(t, argv)
}

func TestNewWithOptionsUsesExplicitInputs(t *testing.T) {
	rt := core.NewRuntime()
	rt.Init()

	NewWithOptions("test", &bytecode.Bytecode{}, rt, Options{
		Args:     []string{"emerald", "script.rb"},
		LoadPath: []string{"/app/lib"},
	})

	assertStringArray(t, rt.Heap.GetGlobalVariableString("$LOAD_PATH"), "/app/lib")
	assertStringArray(t, rt.Heap.GetGlobalVariableString("$*"), "emerald", "script.rb")

	argv, err := getConst(rt.MainObject, "ARGV", rt)
	if err != nil {
		t.Fatal(err)
	}
	assertStringArray(t, argv, "emerald", "script.rb")
}

func assertStringArray(t *testing.T, value object.EmeraldValue, expected ...string) {
	t.Helper()

	array, ok := value.Heap.(*core.ArrayInstance)
	if !ok {
		t.Fatalf("expected Array, got %T", value.Heap)
	}
	if len(array.Value) != len(expected) {
		t.Fatalf("expected %d values, got %d", len(expected), len(array.Value))
	}

	for i, expectedValue := range expected {
		actual, ok := array.Value[i].Heap.(*core.StringInstance)
		if !ok {
			t.Fatalf("value %d: expected String, got %T", i, array.Value[i].Heap)
		}
		if actual.Value != expectedValue {
			t.Fatalf("value %d: expected %q, got %q", i, expectedValue, actual.Value)
		}
	}
}
