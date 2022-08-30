package object

import (
	"fmt"
	"reflect"
)

type BaseEmeraldValue struct {
	builtInMethodSet     BuiltInMethodSet
	definedMethodSet     DefinedMethodSet
	instanceVariables    map[string]EmeraldValue
	includedModules      []EmeraldValue
	parentNamespace      EmeraldValue
	namespaceDefinitions map[string]EmeraldValue
}

func (val *BaseEmeraldValue) IncludedModules() []EmeraldValue { return val.includedModules }

func (val *BaseEmeraldValue) Include(mod EmeraldValue) {
	val.includedModules = append(val.includedModules, mod)
}

func (val *BaseEmeraldValue) NamespaceDefinitions() map[string]EmeraldValue {
	if val.namespaceDefinitions == nil {
		val.namespaceDefinitions = map[string]EmeraldValue{}
	}

	return val.namespaceDefinitions
}

func (val *BaseEmeraldValue) InstanceVariables() map[string]EmeraldValue {
	if val.instanceVariables == nil {
		val.instanceVariables = map[string]EmeraldValue{}
	}

	return val.instanceVariables
}

func (val *BaseEmeraldValue) BuiltInMethodSet() BuiltInMethodSet {
	if val.builtInMethodSet == nil {
		val.builtInMethodSet = BuiltInMethodSet{}
	}

	return val.builtInMethodSet
}

func (val *BaseEmeraldValue) DefinedMethodSet() DefinedMethodSet {
	if val.definedMethodSet == nil {
		val.definedMethodSet = DefinedMethodSet{}
	}

	return val.definedMethodSet
}

func (val *BaseEmeraldValue) DefineMethod(block EmeraldValue, args ...EmeraldValue) {
	name := args[0].Inspect()[1:]

	val.DefinedMethodSet()[name] = block.(*ClosedBlock)
}

func (val *BaseEmeraldValue) Methods(self EmeraldValue) []string {
	methods := []string{}

	for key := range val.BuiltInMethodSet() {
		methods = append(methods, key)
	}

	for key := range val.DefinedMethodSet() {
		methods = append(methods, key)
	}

	for _, mod := range val.IncludedModules() {
		methods = append(methods, mod.Methods(mod)...)
	}

	super := self.Super()
	reflected := reflect.ValueOf(super)

	if super != nil && reflected.IsValid() && !reflected.IsNil() {
		methods = append(methods, super.Methods(super)...)
	}

	return methods
}

func (val *BaseEmeraldValue) RespondsTo(name string, self EmeraldValue) bool {
	_, err := val.ExtractMethod(name, self, self)

	return err == nil
}

func (val *BaseEmeraldValue) SEND(
	ctx *Context,
	yield YieldFunc,
	name string,
	self EmeraldValue,
	block EmeraldValue,
	args ...EmeraldValue,
) (EmeraldValue, error) {
	method, err := self.ExtractMethod(name, self.Class(), self)
	if err != nil {
		return nil, err
	}

	switch method := method.(type) {
	case *ClosedBlock:
		return yield(method, args...), nil
	case *WrappedBuiltInMethod:
		return method.Method(ctx, self, block, yield, args...), nil
	}

	return nil, nil
}

func (val *BaseEmeraldValue) ExtractMethod(name string, extractFrom EmeraldValue, target EmeraldValue) (EmeraldValue, error) {
	if val == nil {
		return nil, fmt.Errorf("invalid method call %s on %#v", name, target)
	}

	for _, ancestor := range extractFrom.Ancestors() {
		if method, ok := ancestor.DefinedMethodSet()[name]; ok {
			return method, nil
		}

		if method, ok := ancestor.BuiltInMethodSet()[name]; ok {
			return &WrappedBuiltInMethod{Method: method}, nil
		}
	}

	return nil, fmt.Errorf("undefined method %s for %s", name, target.Inspect())
}

func (val *BaseEmeraldValue) InstanceVariableGet(name string, extractFrom EmeraldValue, target EmeraldValue) EmeraldValue {
	value, ok := val.InstanceVariables()[name]
	if ok {
		return value
	}

	superClass := extractFrom.Super()
	reflected := reflect.ValueOf(superClass)
	if superClass != nil && reflected.IsValid() && !reflected.IsNil() {
		return superClass.InstanceVariableGet(name, superClass, target)
	}

	return nil
}

func (val *BaseEmeraldValue) InstanceVariableSet(name string, value EmeraldValue) {
	val.InstanceVariables()[name] = value
}

func (val *BaseEmeraldValue) NamespaceDefinitionSet(name string, value EmeraldValue) {
	val.NamespaceDefinitions()[name] = value
}

func (val *BaseEmeraldValue) NamespaceDefinitionGet(name string) EmeraldValue {
	value := val.NamespaceDefinitions()[name]

	if value != nil {
		return value
	}

	if val.parentNamespace != nil {
		return val.parentNamespace.NamespaceDefinitionGet(name)
	}

	return nil
}

func (val *BaseEmeraldValue) ParentNamespace() EmeraldValue {
	return val.parentNamespace
}

func (val *BaseEmeraldValue) SetParentNamespace(parent EmeraldValue) {
	val.parentNamespace = parent
}

func (val *BaseEmeraldValue) ResetForSpec() {
	val.definedMethodSet = DefinedMethodSet{}
	val.instanceVariables = nil
}
