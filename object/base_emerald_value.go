package object

import (
	"fmt"
)

type BaseEmeraldValue struct {
	staticBuiltInMethodSet  BuiltInMethodSet
	staticDefinedMethodSet  DefinedMethodSet
	staticInstanceVariables map[string]EmeraldValue
	builtInMethodSet        BuiltInMethodSet
	definedMethodSet        DefinedMethodSet
	instanceVariables       map[string]EmeraldValue
}

func (val *BaseEmeraldValue) StaticBuiltInMethodSet() BuiltInMethodSet {
	if val.staticBuiltInMethodSet == nil {
		val.staticBuiltInMethodSet = BuiltInMethodSet{}
	}

	return val.staticBuiltInMethodSet
}

func (val *BaseEmeraldValue) StaticDefinedMethodSet() DefinedMethodSet {
	if val.staticDefinedMethodSet == nil {
		val.staticDefinedMethodSet = DefinedMethodSet{}
	}

	return val.staticDefinedMethodSet
}

func (val *BaseEmeraldValue) StaticInstanceVariables() map[string]EmeraldValue {
	if val.staticInstanceVariables == nil {
		val.staticInstanceVariables = map[string]EmeraldValue{}
	}

	return val.staticInstanceVariables
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

func (val *BaseEmeraldValue) DefineMethod(isStatic bool, block EmeraldValue, args ...EmeraldValue) {
	name := args[0].Inspect()[1:]

	var set DefinedMethodSet

	if isStatic {
		set = val.StaticDefinedMethodSet()
	} else {
		set = val.DefinedMethodSet()
	}

	set[name] = block.(*ClosedBlock)
}

func (val *BaseEmeraldValue) RespondsTo(name string, target EmeraldValue) bool {
	_, err := val.ExtractMethod(name, target, target)

	return err == nil
}

func (val *BaseEmeraldValue) SEND(
	yield YieldFunc,
	name string,
	target EmeraldValue,
	block *Block,
	args ...EmeraldValue,
) (EmeraldValue, error) {
	method, err := val.ExtractMethod(name, target, target)
	if err != nil {
		return nil, err
	}

	switch method := method.(type) {
	case *WrappedBuiltInMethod:
		return method.Method(target, block, yield, args...), nil
	}

	return nil, nil
}

func (val *BaseEmeraldValue) ExtractMethod(name string, extractFrom EmeraldValue, target EmeraldValue) (EmeraldValue, error) {
	if val == nil {
		return nil, fmt.Errorf("invalid method call %s on %#v", name, target)
	}

	if _, ok := target.(*Class); ok {
		return val.extractStaticMethod(name, extractFrom, target)
	} else {
		return val.extractInstanceMethod(name, extractFrom, target)
	}
}

func (val *BaseEmeraldValue) extractStaticMethod(name string, extractFrom EmeraldValue, target EmeraldValue) (EmeraldValue, error) {
	if method, ok := val.StaticDefinedMethodSet()[name]; ok {
		return method, nil
	}

	if method, ok := val.StaticBuiltInMethodSet()[name]; ok {
		return &WrappedBuiltInMethod{Method: method}, nil
	}

	superClass := extractFrom.ParentClass().(*Class)

	if superClass != nil {
		super, err := superClass.extractStaticMethod(name, superClass, target)

		if err != nil {
			return nil, fmt.Errorf("undefined method %s for %s:Class", name, target.Inspect())
		}

		return super, nil
	}

	return nil, fmt.Errorf("undefined method %s for %s:Class", name, target.Inspect())
}

func (val *BaseEmeraldValue) extractInstanceMethod(name string, extractFrom EmeraldValue, target EmeraldValue) (EmeraldValue, error) {
	if method, ok := val.DefinedMethodSet()[name]; ok {
		return method, nil
	}

	if method, ok := val.BuiltInMethodSet()[name]; ok {
		return &WrappedBuiltInMethod{Method: method}, nil
	}

	superClass := extractFrom.ParentClass().(*Class)

	if superClass != nil {
		super, err := superClass.extractInstanceMethod(name, superClass, target)

		if err == nil {
			return super, nil
		}
	}

	return nil, fmt.Errorf("undefined method %s for %s:%s", name, target.Inspect(), target.ParentClass().(*Class).Name)
}

func (val *BaseEmeraldValue) InstanceVariableGet(isStatic bool, name string, extractFrom EmeraldValue, target EmeraldValue) EmeraldValue {
	set := val.getInstanceVariables(isStatic)

	value, ok := set[name]
	if ok {
		return value
	}

	superClass := extractFrom.ParentClass().(*Class)

	if superClass != nil {
		return superClass.InstanceVariableGet(isStatic, name, superClass, target)
	}

	return nil
}

func (val *BaseEmeraldValue) InstanceVariableSet(isStatic bool, name string, value EmeraldValue) {
	set := val.getInstanceVariables(isStatic)

	set[name] = value
}

func (val *BaseEmeraldValue) getInstanceVariables(isStatic bool) map[string]EmeraldValue {
	if isStatic {
		return val.StaticInstanceVariables()
	} else {
		return val.InstanceVariables()
	}
}

func (val *BaseEmeraldValue) ResetDefinedMethodSetForSpec() {
	val.definedMethodSet = DefinedMethodSet{}
	val.staticDefinedMethodSet = DefinedMethodSet{}
}
