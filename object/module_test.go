package object

import (
	"testing"
)

func TestModule(t *testing.T) {
	moduleClass := &Class{Name: "Module", BaseEmeraldValue: &BaseEmeraldValue{}}
	myMod := NewModule("MyMod", NewHeapObject(moduleClass), BuiltInMethodSet{}, BuiltInMethodSet{})

	if myMod.Type() != MODULE_VALUE {
		t.Errorf("expected MODULE_VALUE, got %s", myMod.Type().String())
	}

	if myMod.Inspect() != "MyMod" {
		t.Errorf("expected MyMod, got %s", myMod.Inspect())
	}

	if myMod.Class().Heap != moduleClass {
		t.Error("module class not set correctly")
	}

	if !myMod.Super().IsNil() {
		t.Error("module should not have a super")
	}

	ancestors := myMod.Ancestors()
	if len(ancestors) != 1 || ancestors[0].Heap != myMod {
		t.Errorf("unexpected ancestors: %v", ancestors)
	}

	if myMod.SingletonClass().IsNil() {
		t.Error("module should have a singleton class")
	}

	if myMod.HashKey() != "MyMod" {
		t.Error("module hash key should be its name")
	}
}

func TestModule_NewModuleWithStaticMethods(t *testing.T) {
	NewModule("MyMod", EmeraldValue{}, nil, BuiltInMethodSet{"foo": nil})
}

func TestNewModule_WithModules(t *testing.T) {
	mod1 := &Module{Name: "Mod1", BaseEmeraldValue: &BaseEmeraldValue{}}
	mod2 := NewModule("Mod2", EmeraldValue{}, nil, nil, NewHeapObject(mod1))
	if len(mod2.IncludedModules()) != 1 {
		t.Error("module not included")
	}
}

func TestModule_ClassWithSingleton(t *testing.T) {
	myMod := &Module{BaseEmeraldValue: &BaseEmeraldValue{}}
	singleton := myMod.SingletonClass()
	if myMod.Class() != singleton { // Both are EmeraldValue
		t.Error("module class should be its singleton")
	}
}
