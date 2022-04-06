package object

import (
	"emerald/ast"
	"fmt"
)

type (
	BuiltInMethod        func(target EmeraldValue, block *Block, args ...EmeraldValue) EmeraldValue
	WrappedBuiltInMethod struct {
		*BaseEmeraldValue
		Method BuiltInMethod
	}
	BuiltInMethodSet map[string]BuiltInMethod

	DefinedMethodSet map[string]*Block

	EmeraldValueType int
	EmeraldValue     interface {
		Type() EmeraldValueType
		Inspect() string
		ParentClass() EmeraldValue
		DefineMethod(block *Block, args ...EmeraldValue) EmeraldValue
		ExtractMethod(name string, target EmeraldValue) (EmeraldValue, EmeraldValue)
		RespondsTo(name string, target EmeraldValue) bool
		SEND(
			eval func(definitionContext EmeraldValue, node ast.Node, env Environment) EmeraldValue,
			env Environment,
			name string,
			target EmeraldValue,
			block *Block,
			args ...EmeraldValue,
		) EmeraldValue
	}

	BaseEmeraldValue struct {
		builtInMethodSet BuiltInMethodSet
		definedMethodSet DefinedMethodSet
	}
)

const (
	CLASS_VALUE EmeraldValueType = iota
	INSTANCE_VALUE
	BLOCK_VALUE
	RETURN_VALUE
)

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

func (val *BaseEmeraldValue) DefineMethod(block *Block, args ...EmeraldValue) EmeraldValue {
	name := args[0].Inspect()

	val.DefinedMethodSet()[name] = block

	return NewSymbol(name)
}

func (val *BaseEmeraldValue) RespondsTo(name string, target EmeraldValue) bool {
	_, err := val.ExtractMethod(name, target)

	return err == nil
}

func (val *BaseEmeraldValue) SEND(
	eval func(definitionContext EmeraldValue, node ast.Node, env Environment) EmeraldValue,
	env Environment,
	name string,
	target EmeraldValue,
	block *Block,
	args ...EmeraldValue,
) EmeraldValue {
	method, err := val.ExtractMethod(name, target)
	if err != nil {
		return err
	}

	switch method := method.(type) {
	case *WrappedBuiltInMethod:
		return method.Method(target, block, args...)
	case *Block:
		evaluated := eval(target, method.Body, extendBlockEnv(method, args))
		return unwrapReturnValue(evaluated)
	}

	return nil
}

func (val *BaseEmeraldValue) ExtractMethod(name string, target EmeraldValue) (EmeraldValue, EmeraldValue) {
	if method, ok := val.DefinedMethodSet()[name]; ok {
		return method, nil
	}

	if method, ok := val.BuiltInMethodSet()[name]; ok {
		return &WrappedBuiltInMethod{Method: method}, nil
	}

	superClass := target.ParentClass().(*Class)

	if superClass != nil {
		super, err := superClass.ExtractMethod(name, superClass)

		if err != nil {
			return nil, NewStandardError(
				fmt.Sprintf("undefined method %s for %s:%s", name, target.Inspect(), superClass.Name),
			)
		}

		return super, nil
	}

	return nil, NewStandardError(
		fmt.Sprintf("undefined method %s for %s:%s", name, target.Inspect(), target.(*Class).Name),
	)
}

func (method *WrappedBuiltInMethod) Inspect() string           { return "obscure Go code" }
func (method *WrappedBuiltInMethod) Type() EmeraldValueType    { return BLOCK_VALUE }
func (method *WrappedBuiltInMethod) ParentClass() EmeraldValue { return nil }

func extendBlockEnv(
	fn *Block,
	args []EmeraldValue,
) Environment {
	env := NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.(*ast.IdentifierExpression).Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj EmeraldValue) EmeraldValue {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func nativeBoolToBooleanObject(input bool) EmeraldValue {
	if input {
		return TRUE
	}
	return FALSE
}
