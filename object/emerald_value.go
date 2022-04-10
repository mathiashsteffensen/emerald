package object

import (
	"emerald/ast"
)

type (
	// BuiltInMethod - The type signature of an Emerald method defined in Go compiler
	BuiltInMethod func(target EmeraldValue, block *Block, yield YieldFunc, args ...EmeraldValue) EmeraldValue

	// WrappedBuiltInMethod -  Wraps a built-in method so that it conforms to the EmeraldValue interface
	WrappedBuiltInMethod struct {
		*BaseEmeraldValue
		Method BuiltInMethod
	}

	// BuiltInMethodSet - Stores an objects built-in method set
	BuiltInMethodSet map[string]BuiltInMethod

	// DefinedMethodSet - Stores an objects methods defined by the program
	DefinedMethodSet map[string]*Block

	EmeraldValueType int
	// EmeraldValue - All Emerald objects must implement this interface
	EmeraldValue interface {
		Type() EmeraldValueType
		Inspect() string
		ParentClass() EmeraldValue
		DefineMethod(isStatic bool, block *Block, args ...EmeraldValue) EmeraldValue
		ExtractMethod(name string, extractFrom EmeraldValue, target EmeraldValue) (EmeraldValue, EmeraldValue)
		RespondsTo(name string, target EmeraldValue) bool
		SEND(
			eval func(executionContext ExecutionContext, node ast.Node, env Environment) EmeraldValue,
			yield YieldFunc,
			name string,
			target EmeraldValue,
			block *Block,
			args ...EmeraldValue,
		) EmeraldValue
		InstanceVariableGet(isStatic bool, name string, extractFrom EmeraldValue, target EmeraldValue) EmeraldValue
		InstanceVariableSet(isStatic bool, name string, value EmeraldValue)
	}
)

const (
	CLASS_VALUE EmeraldValueType = iota
	INSTANCE_VALUE
	BLOCK_VALUE
	RETURN_VALUE
)

func (method *WrappedBuiltInMethod) Inspect() string           { return "obscure Go compiler" }
func (method *WrappedBuiltInMethod) Type() EmeraldValueType    { return BLOCK_VALUE }
func (method *WrappedBuiltInMethod) ParentClass() EmeraldValue { return nil }

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
