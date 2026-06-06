package object

import (
	"testing"
)

func TestBaseEmeraldValue_IncludedModules(t *testing.T) {
	mod1 := &Module{BaseEmeraldValue: &BaseEmeraldValue{}, Name: "Mod1"}
	mod2 := &Module{BaseEmeraldValue: &BaseEmeraldValue{}, Name: "Mod2"}
	mod3 := &Module{BaseEmeraldValue: &BaseEmeraldValue{}, Name: "Mod3"}

	mod1.Include(NewHeapObject(mod2))
	mod2.Include(NewHeapObject(mod3))

	base := &Class{BaseEmeraldValue: &BaseEmeraldValue{}, Name: "Base"}
	base.Include(NewHeapObject(mod1))

	included := base.IncludedModules()
	if len(included) != 3 {
		t.Errorf("expected 3 included modules, got %d", len(included))
	}

	if included[0] != NewHeapObject(mod1) || included[1] != NewHeapObject(mod2) || included[2] != NewHeapObject(mod3) {
		t.Errorf("unexpected order of included modules")
	}
}

func TestBaseEmeraldValue_NamespaceDefinitions(t *testing.T) {
	parent := &Class{BaseEmeraldValue: &BaseEmeraldValue{}, Name: "Parent"}
	child := &Class{BaseEmeraldValue: &BaseEmeraldValue{}, Name: "Child"}
	child.SetParentNamespace(NewHeapObject(parent))

	val := &Class{Name: "MyClass"}
	parent.NamespaceDefinitionSet("MyClass", NewHeapObject(val))

	if child.NamespaceDefinitionGet("MyClass") != NewHeapObject(val) {
		t.Error("failed to get namespace definition from parent")
	}

	if !child.NamespaceDefinitionGet("NonExistent").IsNil() {
		t.Error("expected nil for non-existent namespace definition")
	}

	if child.ParentNamespace() != NewHeapObject(parent) {
		t.Error("parent namespace not set correctly")
	}
}

func TestBaseEmeraldValue_InstanceVariables(t *testing.T) {
	class := &Class{Name: "MyClass", BaseEmeraldValue: &BaseEmeraldValue{}}
	subClass := &Class{Name: "SubClass", BaseEmeraldValue: &BaseEmeraldValue{}, super: NewHeapObject(class)}
	instance := &Instance{BaseEmeraldValue: &BaseEmeraldValue{}, baseClass: NewHeapObject(subClass)}

	val := &Class{Name: "SomeValue"}
	instance.InstanceVariableSet("@var", NewHeapObject(val))

	if instance.InstanceVariableGet("@var", NewHeapObject(instance), NewHeapObject(instance)) != NewHeapObject(val) {
		t.Errorf("failed to get instance variable from instance, got %v", instance.InstanceVariableGet("@var", NewHeapObject(instance), NewHeapObject(instance)))
	}

	subClass.InstanceVariableSet("@sub_var", NewHeapObject(val))
	if subClass.InstanceVariableGet("@sub_var", NewHeapObject(subClass), NewHeapObject(subClass)) != NewHeapObject(val) {
		t.Errorf("failed to get instance variable from sub class, got %v", subClass.InstanceVariableGet("@sub_var", NewHeapObject(subClass), NewHeapObject(subClass)))
	}

	if !subClass.InstanceVariableGet("@non_existent", NewHeapObject(subClass), NewHeapObject(subClass)).IsNil() {
		t.Error("expected nil for non-existent instance variable")
	}
}

func TestBaseEmeraldValue_InstanceVariableInheritance(t *testing.T) {
	class := &Class{Name: "MyClass", BaseEmeraldValue: &BaseEmeraldValue{}}
	subClass := &Class{Name: "SubClass", BaseEmeraldValue: &BaseEmeraldValue{}, super: NewHeapObject(class)}

	val := &Class{Name: "SomeValue"}
	class.InstanceVariableSet("@class_var", NewHeapObject(val))

	if subClass.InstanceVariableGet("@class_var", NewHeapObject(subClass), NewHeapObject(subClass)) != NewHeapObject(val) {
		t.Error("failed to get instance variable from super class")
	}
}

func TestBaseEmeraldValue_Methods(t *testing.T) {
	base := &BaseEmeraldValue{}

	method1 := &WrappedBuiltInMethod{}
	base.BuiltInMethodSet()["method1"] = method1

	block := &ClosedBlock{Block: &Block{}}
	// DefineMethod expects a name in args[0] (wrapped in EmeraldValue)
	name := NewHeapObject(&Class{Name: ":method2"}) // Inspect() will return "method2", [1:] will return "method2"
	// Wait, let's check DefineMethod implementation:
	// func (val *BaseEmeraldValue) DefineMethod(block EmeraldValue, args ...EmeraldValue) {
	// 	name := args[0].Inspect()[1:]
	// 	val.DefinedMethodSet()[name] = block.(*ClosedBlock)
	// }
	base.DefineMethod(NewHeapObject(block), name)

	methods := base.Methods()
	if len(methods) != 2 {
		t.Errorf("expected 2 methods, got %d", len(methods))
	}

	found1 := false
	found2 := false
	for _, m := range methods {
		if m == "method1" {
			found1 = true
		}
		if m == "method2" {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Error("did not find all methods")
	}
}

func TestBaseEmeraldValue_ExtractMethod(t *testing.T) {
	class := &Class{
		Name: "MyClass",
		BaseEmeraldValue: &BaseEmeraldValue{
			builtInMethodSet: BuiltInMethodSet{
				"my_method": &WrappedBuiltInMethod{Visibility: PUBLIC},
			},
		},
	}
	instance := &Instance{BaseEmeraldValue: &BaseEmeraldValue{}, baseClass: NewHeapObject(class)}

	method, visibility, isSelf, err := instance.ExtractMethod("my_method", NewHeapObject(instance), NewHeapObject(instance))
	if err != nil {
		t.Fatalf("failed to extract method: %s", err)
	}

	if visibility != PUBLIC {
		t.Errorf("expected PUBLIC visibility, got %s", visibility)
	}

	if isSelf {
		t.Error("expected isSelf to be false (method is in class, not instance)")
	}

	if method.IsNil() {
		t.Error("extracted method is nil")
	}

	_, _, _, err = instance.ExtractMethod("non_existent", NewHeapObject(instance), NewHeapObject(instance))
	if err == nil {
		t.Error("expected error when extracting non-existent method")
	}
}

func TestBaseEmeraldValue_RespondsTo(t *testing.T) {
	class := &Class{
		Name: "MyClass",
		BaseEmeraldValue: &BaseEmeraldValue{
			builtInMethodSet: BuiltInMethodSet{
				"public_method":  &WrappedBuiltInMethod{Visibility: PUBLIC},
				"private_method": &WrappedBuiltInMethod{Visibility: PRIVATE},
			},
		},
	}
	instance := &Instance{BaseEmeraldValue: &BaseEmeraldValue{}, baseClass: NewHeapObject(class)}

	if !NewHeapObject(instance).RespondsTo("public_method", NewHeapObject(instance)) {
		t.Error("expected RespondsTo public_method to be true")
	}

	if NewHeapObject(instance).RespondsTo("private_method", NewHeapObject(instance)) {
		t.Error("expected RespondsTo private_method to be false")
	}

	if NewHeapObject(instance).RespondsTo("non_existent", NewHeapObject(instance)) {
		t.Error("expected RespondsTo non_existent to be false")
	}
}

func TestBaseEmeraldValue_ResetForSpec(t *testing.T) {
	base := &Class{BaseEmeraldValue: &BaseEmeraldValue{}, Name: "Base"}
	base.InstanceVariableSet("@var", NewHeapObject(&Class{Name: "Val"}))
	base.ResetForSpec()

	if len(base.InstanceVariables()) != 0 {
		t.Error("ResetForSpec did not clear instance variables")
	}

	base.ResetForSpec() // Should handle nil instanceVariables
}

func TestBaseEmeraldValue_DefinedMethodSet(t *testing.T) {
	base := &BaseEmeraldValue{}
	base.DefinedMethodSet()
}
