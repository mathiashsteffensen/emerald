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

func (val *BaseEmeraldValue) IncludedModules() []EmeraldValue {
	direct := val.includedModules

	for _, module := range direct {
		direct = append(direct, module.IncludedModules()...)
	}

	return direct
}

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

func (val *BaseEmeraldValue) Methods() []string {
	methods := []string{}

	for key := range val.BuiltInMethodSet() {
		methods = append(methods, key)
	}

	for key := range val.DefinedMethodSet() {
		methods = append(methods, key)
	}

	return methods
}

func (val *BaseEmeraldValue) RespondsTo(name string, self EmeraldValue) bool {
	_, err := val.ExtractMethod(name, self, self)

	return err == nil
}

var EvalBlock func(block *ClosedBlock, args ...EmeraldValue) EmeraldValue

func (val *BaseEmeraldValue) SEND(
	ctx *Context,
	name string,
	args ...EmeraldValue,
) EmeraldValue {
	method, err := ctx.Self.ExtractMethod(name, ctx.Self.Class(), ctx.Self)
	if err != nil {
		panic(err)
	}

	switch method := method.(type) {
	case *ClosedBlock:
		return EvalBlock(method, args...)
	case *WrappedBuiltInMethod:
		return method.Method(ctx, args...)
	default:
		panic("This is a bug in the Emerald VM, no idea how the fuck we got here, my b")
	}
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
