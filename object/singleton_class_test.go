package object

import (
	"testing"
)

func TestSingletonClass(t *testing.T) {
	super := &Class{Name: "Object", BaseEmeraldValue: &BaseEmeraldValue{}}
	instance := &Instance{BaseEmeraldValue: &BaseEmeraldValue{}, baseClass: NewHeapObject(super)}
	singleton := NewSingletonClass(NewHeapObject(instance), NewHeapObject(super), &BaseEmeraldValue{})

	if singleton.Type() != STATIC_CLASS_VALUE {
		t.Errorf("expected STATIC_CLASS_VALUE, got %s", singleton.Type().String())
	}

	if singleton.Super().Heap != super {
		t.Error("singleton super not set correctly")
	}

	if !singleton.SingletonClass().IsNil() {
		t.Error("singleton of singleton should be nil")
	}

	if singleton.Class().Heap != super {
		t.Errorf("expected class %v, got %v", super, singleton.Class())
	}

	singleton.SetSuper(EmeraldValue{})
	if !singleton.Super().IsNil() {
		t.Error("SetSuper failed")
	}

	if len(singleton.Ancestors()) == 0 { // This might now fail because super is nil
		t.Error("singleton should have ancestors")
	}

	if singleton.HashKey() == "" {
		t.Error("singleton hash key should not be empty")
	}

	if singleton.Inspect() == "" {
		t.Error("singleton inspect should not be empty")
	}
}
