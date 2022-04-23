package object

import (
	"fmt"
	"reflect"
)

type BaseEmeraldValue struct {
	builtInMethodSet  BuiltInMethodSet
	definedMethodSet  DefinedMethodSet
	instanceVariables map[string]EmeraldValue
	includedModules   []EmeraldValue
}

func (val *BaseEmeraldValue) IncludedModules() []EmeraldValue { return val.includedModules }

func (val *BaseEmeraldValue) Include(mod EmeraldValue) {
	val.includedModules = append(val.includedModules, mod)
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

func (val *BaseEmeraldValue) Methods(target EmeraldValue) []string {
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

	super := target.ParentClass()
	reflected := reflect.ValueOf(super)

	if super != nil && reflected.IsValid() && !reflected.IsNil() {
		methods = append(methods, super.Methods(super)...)
	}

	return methods
}

func (val *BaseEmeraldValue) RespondsTo(name string, target EmeraldValue) bool {
	_, err := val.ExtractMethod(name, target, target)

	return err == nil
}

func (val *BaseEmeraldValue) SEND(
	ctx *Context,
	yield YieldFunc,
	name string,
	target EmeraldValue,
	block *ClosedBlock,
	args ...EmeraldValue,
) (EmeraldValue, error) {
	method, err := val.ExtractMethod(name, target, target)
	if err != nil {
		return nil, err
	}

	switch method := method.(type) {
	case *WrappedBuiltInMethod:
		return method.Method(ctx, target, block, yield, args...), nil
	}

	return nil, nil
}

func (val *BaseEmeraldValue) ExtractMethod(name string, extractFrom EmeraldValue, target EmeraldValue) (EmeraldValue, error) {
	if val == nil {
		return nil, fmt.Errorf("invalid method call %s on %#v", name, target)
	}

	if targetInstance, ok := target.(*Instance); ok {
		method, ok := targetInstance.BuiltInSingletonMethods[name]
		if ok {
			return &WrappedBuiltInMethod{Method: method}, nil
		}
	}

	if method, ok := val.DefinedMethodSet()[name]; ok {
		return method, nil
	}

	if method, ok := val.BuiltInMethodSet()[name]; ok {
		return &WrappedBuiltInMethod{Method: method}, nil
	}

	for _, mod := range val.IncludedModules() {
		if method, err := mod.ExtractMethod(name, mod, target); err == nil {
			return method, nil
		}
	}

	super := extractFrom.ParentClass()
	reflected := reflect.ValueOf(super)

	if super != nil && reflected.IsValid() && !reflected.IsNil() {
		superMethod, err := super.ExtractMethod(name, super, target)

		if err == nil {
			return superMethod, nil
		}
	}

	var parentName string
	super = target.ParentClass()
	reflected = reflect.ValueOf(super)

	if super != nil && reflected.IsValid() && !reflected.IsNil() {
		switch super := super.(type) {
		case *Class:
			parentName = super.Name
		case *StaticClass:
			parentName = super.Name
		}
	}

	return nil, fmt.Errorf("undefined method %s for %s:%s", name, target.Inspect(), parentName)
}

func (val *BaseEmeraldValue) InstanceVariableGet(name string, extractFrom EmeraldValue, target EmeraldValue) EmeraldValue {
	value, ok := val.instanceVariables[name]
	if ok {
		return value
	}

	superClass := extractFrom.ParentClass()
	reflected := reflect.ValueOf(superClass)
	if superClass != nil && reflected.IsValid() && !reflected.IsNil() {
		return superClass.InstanceVariableGet(name, superClass, target)
	}

	return nil
}

func (val *BaseEmeraldValue) InstanceVariableSet(name string, value EmeraldValue) {
	val.InstanceVariables()[name] = value
}

func (val *BaseEmeraldValue) ResetDefinedMethodSetForSpec() {
	val.definedMethodSet = DefinedMethodSet{}
}
