package object

import (
	"fmt"
)

type BaseEmeraldValue struct {
	builtInMethodSet     BuiltInMethodSet
	definedMethodSet     DefinedMethodSet
	instanceVariables    map[string]EmeraldValue
	includedModules      []EmeraldValue
	parentNamespace      EmeraldValue
	namespaceDefinitions map[string]EmeraldValue
}

func (val *BaseEmeraldValue) SetType(t EmeraldValueType) {}
func (val *BaseEmeraldValue) SetHeap(h HeapObject)       {}

func (val *BaseEmeraldValue) IncludedModules() []EmeraldValue {
	included := []EmeraldValue{}
	collectIncludedModules(val.includedModules, map[HeapObject]struct{}{}, &included)
	return included
}

func collectIncludedModules(
	modules []EmeraldValue,
	visited map[HeapObject]struct{},
	included *[]EmeraldValue,
) {
	direct := make([]EmeraldValue, 0, len(modules))

	for _, module := range modules {
		if module.Heap == nil {
			continue
		}
		if _, seen := visited[module.Heap]; seen {
			continue
		}

		visited[module.Heap] = struct{}{}
		*included = append(*included, module)
		direct = append(direct, module)
	}

	for _, module := range direct {
		if provider, ok := module.Heap.(baseEmeraldValueProvider); ok {
			collectIncludedModules(provider.baseEmeraldValue().includedModules, visited, included)
		}
	}
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

	val.DefinedMethodSet()[name] = block.Heap.(*ClosedBlock)
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
	_, visibility, _, err := val.ExtractMethod(name, self, self)

	return err == nil && visibility == PUBLIC
}

func (val *BaseEmeraldValue) ExtractMethod(name string, extractFrom EmeraldValue, target EmeraldValue) (EmeraldValue, MethodVisibility, bool, error) {
	for _, ancestor := range extractFrom.Ancestors() {
		isSelf := ancestor == extractFrom

		if method, ok := ancestor.DefinedMethodSet()[name]; ok {
			return NewHeapObject(method), method.Visibility, isSelf, nil
		}

		if method, ok := ancestor.BuiltInMethodSet()[name]; ok {
			return NewHeapObject(method), method.Visibility, isSelf, nil
		}
	}

	return EmeraldValue{}, PUBLIC, false, fmt.Errorf("undefined method %s for %s", name, target.Inspect())
}

func (val *BaseEmeraldValue) InstanceVariableGet(name string, extractFrom EmeraldValue, target EmeraldValue) EmeraldValue {
	value, ok := val.InstanceVariables()[name]
	if ok {
		return value
	}

	superClass := extractFrom.Super()
	if !superClass.IsNil() {
		return superClass.InstanceVariableGet(name, superClass, target)
	}

	return EmeraldValue{}
}

func (val *BaseEmeraldValue) InstanceVariableSet(name string, value EmeraldValue) {
	val.InstanceVariables()[name] = value
}

func (val *BaseEmeraldValue) NamespaceDefinitionSet(name string, value EmeraldValue) {
	val.NamespaceDefinitions()[name] = value
}

func (val *BaseEmeraldValue) NamespaceDefinitionGet(name string) EmeraldValue {
	visited := map[*BaseEmeraldValue]struct{}{}
	current := val

	for current != nil {
		if _, seen := visited[current]; seen {
			return EmeraldValue{}
		}
		visited[current] = struct{}{}

		value := current.NamespaceDefinitions()[name]
		if !value.IsNil() {
			return value
		}

		parent := current.parentNamespace
		if parent.IsNil() || parent.Heap == nil {
			return EmeraldValue{}
		}

		provider, ok := parent.Heap.(baseEmeraldValueProvider)
		if !ok {
			return EmeraldValue{}
		}
		current = provider.baseEmeraldValue()
	}

	return EmeraldValue{}
}

type baseEmeraldValueProvider interface {
	baseEmeraldValue() *BaseEmeraldValue
}

func (val *BaseEmeraldValue) baseEmeraldValue() *BaseEmeraldValue {
	return val
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
