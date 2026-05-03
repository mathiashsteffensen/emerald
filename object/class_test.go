package object

import (
	"testing"
)

func TestClass(t *testing.T) {
	objectClass := &Class{Name: "Object", BaseEmeraldValue: &BaseEmeraldValue{}}
	myClass := NewClass("MyClass", objectClass, nil, BuiltInMethodSet{}, BuiltInMethodSet{})

	if myClass.Type() != CLASS_VALUE {
		t.Errorf("expected CLASS_VALUE, got %s", myClass.Type().String())
	}

	if myClass.Inspect() != "MyClass" {
		t.Errorf("expected MyClass, got %s", myClass.Inspect())
	}

	if myClass.Super() != objectClass {
		t.Error("super class not set correctly")
	}

	ancestors := myClass.Ancestors()
	if len(ancestors) != 2 || ancestors[0] != myClass || ancestors[1] != objectClass {
		t.Errorf("unexpected ancestors: %v", ancestors)
	}

	if myClass.HashKey() != "MyClass" {
		t.Errorf("expected MyClass hash key, got %s", myClass.HashKey())
	}

	instance := myClass.New()
	if instance.baseClass != myClass {
		t.Error("instance base class not set correctly")
	}
}

func TestClass_SingletonClass(t *testing.T) {
	myClass := &Class{Name: "MyClass", BaseEmeraldValue: &BaseEmeraldValue{}}
	singleton := myClass.SingletonClass()

	if singleton.Type() != STATIC_CLASS_VALUE {
		t.Errorf("expected STATIC_CLASS_VALUE, got %s", singleton.Type().String())
	}

	if myClass.Class() != singleton {
		t.Error("Class() should return singleton if it exists")
	}

	if singleton.SingletonClass() != nil {
		t.Error("singleton of singleton should be nil")
	}
}

func TestClass_SetSuper(t *testing.T) {
	class := &Class{}
	super := &Class{}
	class.SetSuper(super)
}

func TestNewClass_WithModules(t *testing.T) {
	mod := &Module{Name: "MyMod", BaseEmeraldValue: &BaseEmeraldValue{}}
	class := NewClass("MyClass", nil, nil, nil, nil, mod)
	if len(class.IncludedModules()) != 1 {
		t.Error("module not included")
	}
}
