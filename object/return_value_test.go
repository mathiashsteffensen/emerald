package object

import (
	"testing"
)

func TestReturnValue(t *testing.T) {
	val := &Class{Name: "MyClass"}
	rv := &ReturnValue{Value: NewHeapObject(val)}

	if rv.Type() != RETURN_VALUE {
		t.Errorf("expected RETURN_VALUE, got %s", rv.Type().String())
	}

	if rv.Value.Inspect() != "MyClass" {
		t.Errorf("unexpected value inspect: %s", rv.Value.Inspect())
	}

	if !rv.Class().IsNil() {
		t.Error("return value class should be nil")
	}

	if !rv.Super().IsNil() {
		t.Error("return value super should be nil")
	}

	if len(rv.Ancestors()) != 0 {
		t.Error("return value ancestors should be empty")
	}

	if rv.HashKey() == "" {
		t.Error("return value hash key should not be empty")
	}

	if !rv.SingletonClass().IsNil() {
		t.Error("return value singleton class should be nil")
	}
}
