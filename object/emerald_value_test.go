package object

import (
	"testing"
	"unsafe"
)

func TestEmeraldValueSize(t *testing.T) {
	var value EmeraldValue

	if unsafe.Sizeof(value) != 16 {
		t.Errorf("blank EmeraldValue takes up %d bytes", unsafe.Sizeof(value))
	}
}

func TestEmeraldValueType_String(t *testing.T) {
	tests := []struct {
		t        EmeraldValueType
		expected string
	}{
		{CLASS_VALUE, "Class"},
		{STATIC_CLASS_VALUE, "Static Class"},
		{MODULE_VALUE, "Module"},
		{INSTANCE_VALUE, "Instance"},
		{BLOCK_VALUE, "Block"},
		{RETURN_VALUE, "Return"},
		{EmeraldValueType(99), ""},
	}

	for _, tt := range tests {
		if tt.t.String() != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, tt.t.String())
		}
	}
}

func TestRealClass(t *testing.T) {
	if RealClass(nil) != nil {
		t.Error("RealClass(nil) should be nil")
	}

	class := &Class{Name: "MyClass"}
	if RealClass(class) != nil {
		t.Error("RealClass(class) should be nil if Class() is nil")
	}

	singleton := &SingletonClass{BaseEmeraldValue: &BaseEmeraldValue{}, super: class}
	instance := &Instance{BaseEmeraldValue: &BaseEmeraldValue{}, singleton: singleton}
	if RealClass(instance) != class {
		t.Errorf("expected RealClass(instance) to be class, got %v", RealClass(instance))
	}

	if RealClass(singleton) != class { 
		t.Errorf("expected RealClass(singleton) to be class, got %v", RealClass(singleton))
	}
}

func TestWrappedBuiltInMethod(t *testing.T) {
	method := &WrappedBuiltInMethod{}
	if method.Type() != BLOCK_VALUE {
		t.Errorf("expected BLOCK_VALUE, got %s", method.Type().String())
	}
	if method.Class() != nil {
		t.Error("expected nil Class")
	}
	if method.Super() != nil {
		t.Error("expected nil Super")
	}
	if len(method.Ancestors()) != 0 {
		t.Error("expected empty Ancestors")
	}
	if method.SingletonClass() != nil {
		t.Error("expected nil SingletonClass")
	}
	if method.HashKey() == "" {
		t.Error("expected non-empty HashKey")
	}
}

