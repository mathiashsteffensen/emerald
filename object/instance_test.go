package object

import (
	"strings"
	"testing"
)

func TestInstance(t *testing.T) {
	myClass := &Class{Name: "MyClass", BaseEmeraldValue: &BaseEmeraldValue{}}
	instance := myClass.New()

	if instance.Type() != INSTANCE_VALUE {
		t.Errorf("expected INSTANCE_VALUE, got %s", instance.Type().String())
	}

	if instance.Class().Heap != myClass {
		t.Errorf("expected class MyClass, got %v", instance.Class())
	}

	if !instance.Super().IsNil() {
		t.Error("instance should not have a super")
	}

	if !strings.HasPrefix(instance.Inspect(), "#<MyClass:0x") {
		t.Errorf("unexpected inspect format: %s", instance.Inspect())
	}

	if instance.HashKey() != instance.Inspect() {
		t.Error("HashKey should match Inspect")
	}
}

func TestInstance_SingletonClass(t *testing.T) {
	myClass := &Class{Name: "MyClass", BaseEmeraldValue: &BaseEmeraldValue{}}
	instance := myClass.New()

	singleton := instance.SingletonClass()
	if instance.Class() != singleton { // Both are EmeraldValue
		t.Error("instance class should be its singleton")
	}

	if singleton.Super().Heap != myClass {
		t.Error("singleton super should be the base class")
	}
}

func TestInstance_Ancestors(t *testing.T) {
	objectClass := &Class{Name: "Object", BaseEmeraldValue: &BaseEmeraldValue{}}
	myClass := &Class{Name: "MyClass", BaseEmeraldValue: &BaseEmeraldValue{}, super: NewHeapObject(objectClass)}
	instance := myClass.New()

	ancestors := instance.Ancestors()
	// Ancestors for instance: Class.Ancestors() + instance
	// myClass.Ancestors() = [myClass, objectClass]
	// instance.Ancestors() = [myClass, objectClass, instance]

	if len(ancestors) != 3 {
		t.Fatalf("expected 3 ancestors, got %d", len(ancestors))
	}
}

func TestInstance_Include(t *testing.T) {
	instance := &Instance{BaseEmeraldValue: &BaseEmeraldValue{}}
	mod := &Module{Name: "MyMod", BaseEmeraldValue: &BaseEmeraldValue{}}

	instance.Include(NewHeapObject(mod))

	if len(instance.SingletonClass().IncludedModules()) != 1 {
		t.Error("module not included in singleton class")
	}
}

func TestInstance_DefineMethod(t *testing.T) {
	instance := &Instance{BaseEmeraldValue: &BaseEmeraldValue{}}
	block := &ClosedBlock{Block: &Block{}}
	name := &Class{Name: ":my_method"}

	instance.DefineMethod(NewHeapObject(block), NewHeapObject(name))
}
