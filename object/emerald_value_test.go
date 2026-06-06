package object

import (
	"testing"
	"unsafe"
)

func TestEmeraldValueSize(t *testing.T) {
	var value EmeraldValue

	if unsafe.Sizeof(value) > 32 {
		t.Errorf("EmeraldValue takes up too much space: %d bytes", unsafe.Sizeof(value))
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
	if !RealClass(EmeraldValue{}).IsNil() {
		t.Error("RealClass(nil) should be nil")
	}

	class := &Class{Name: "MyClass"}
	if !RealClass(NewHeapObject(class)).IsNil() {
		t.Error("RealClass(class) should be nil if Class() is nil")
	}

	singleton := &SingletonClass{BaseEmeraldValue: &BaseEmeraldValue{}, super: NewHeapObject(class)}
	instance := &Instance{BaseEmeraldValue: &BaseEmeraldValue{}, singleton: singleton}
	if RealClass(NewHeapObject(instance)).Heap != class {
		t.Errorf("expected RealClass(instance) to be class, got %v", RealClass(NewHeapObject(instance)))
	}

	if RealClass(NewHeapObject(singleton)).Heap != class {
		t.Errorf("expected RealClass(singleton) to be class, got %v", RealClass(NewHeapObject(singleton)))
	}
}

func TestWrappedBuiltInMethod(t *testing.T) {
	method := &WrappedBuiltInMethod{
		BaseEmeraldValue: &BaseEmeraldValue{},
	}
	if method.Type() != BLOCK_VALUE {
		t.Errorf("expected BLOCK_VALUE, got %s", method.Type().String())
	}
	if !method.Class().IsNil() {
		t.Error("expected nil Class")
	}
	if !method.Super().IsNil() {
		t.Error("expected nil Super")
	}
	// Ancestors() on BaseEmeraldValue (which WrappedBuiltInMethod uses via embedding)
	// might return something if it's not overridden.
	if len(method.Ancestors()) != 0 {
		t.Error("expected empty Ancestors")
	}
	if !method.SingletonClass().IsNil() {
		t.Error("expected nil SingletonClass")
	}
	if method.HashKey() == "" {
		t.Error("expected non-empty HashKey")
	}
}
