package object

import (
	"testing"
)

func TestReturnValue(t *testing.T) {
	val := &Class{Name: "MyClass"}
	rv := &ReturnValue{Value: val}

	if rv.Type() != RETURN_VALUE {
		t.Errorf("expected RETURN_VALUE, got %s", rv.Type().String())
	}

	if rv.Inspect() != "return MyClass" {
		t.Errorf("unexpected inspect: %s", rv.Inspect())
	}

	if rv.Class() != nil {
		t.Error("return value class should be nil")
	}

	if rv.Super() != nil {
		t.Error("return value super should be nil")
	}

	if len(rv.Ancestors()) != 0 {
		t.Error("return value ancestors should be empty")
	}

	if rv.HashKey() == "" {
		t.Error("return value hash key should not be empty")
	}

	if rv.SingletonClass() != nil {
		t.Error("return value singleton class should be nil")
	}
}
