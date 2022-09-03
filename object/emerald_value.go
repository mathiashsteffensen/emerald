package object

import "fmt"

type (
	// BuiltInMethod - The type signature of an Emerald method defined in Go compiler
	BuiltInMethod func(ctx *Context, args ...EmeraldValue) EmeraldValue

	// WrappedBuiltInMethod -  Wraps a built-in method so that it conforms to the EmeraldValue interface
	WrappedBuiltInMethod struct {
		*BaseEmeraldValue
		Method BuiltInMethod
	}

	// BuiltInMethodSet - Stores an objects built-in method set
	BuiltInMethodSet map[string]BuiltInMethod

	// DefinedMethodSet - Stores an objects methods defined by the program
	DefinedMethodSet map[string]*ClosedBlock

	EmeraldValueType int
	// EmeraldValue - All Emerald objects must implement this interface
	EmeraldValue interface {
		Type() EmeraldValueType
		Inspect() string
		Class() EmeraldValue
		Super() EmeraldValue
		Ancestors() []EmeraldValue
		IncludedModules() []EmeraldValue
		Include(mod EmeraldValue)
		BuiltInMethodSet() BuiltInMethodSet
		DefinedMethodSet() DefinedMethodSet
		DefineMethod(block EmeraldValue, args ...EmeraldValue)
		ExtractMethod(name string, extractFrom EmeraldValue, self EmeraldValue) (EmeraldValue, error)
		Methods() []string
		RespondsTo(name string, self EmeraldValue) bool
		SEND(
			ctx *Context,
			name string,
			args ...EmeraldValue,
		) EmeraldValue
		InstanceVariableGet(name string, extractFrom EmeraldValue, self EmeraldValue) EmeraldValue
		InstanceVariableSet(name string, value EmeraldValue)
		ParentNamespace() EmeraldValue
		SetParentNamespace(parent EmeraldValue)
		NamespaceDefinitionGet(name string) EmeraldValue
		NamespaceDefinitionSet(name string, value EmeraldValue)
		HashKey() string
	}
)

const (
	_ EmeraldValueType = iota
	CLASS_VALUE
	STATIC_CLASS_VALUE
	MODULE_VALUE
	INSTANCE_VALUE
	BLOCK_VALUE
	RETURN_VALUE
)

func (t EmeraldValueType) String() string {
	switch t {
	case CLASS_VALUE:
		return "Class"
	case STATIC_CLASS_VALUE:
		return "Static Class"
	case MODULE_VALUE:
		return "Module"
	case INSTANCE_VALUE:
		return "Instance"
	case BLOCK_VALUE:
		return "Block"
	case RETURN_VALUE:
		return "Return"
	}

	return ""
}

func (method *WrappedBuiltInMethod) Inspect() string           { return fmt.Sprintf("#<Block:%p>", method) }
func (method *WrappedBuiltInMethod) Type() EmeraldValueType    { return BLOCK_VALUE }
func (method *WrappedBuiltInMethod) Class() EmeraldValue       { return nil }
func (method *WrappedBuiltInMethod) Super() EmeraldValue       { return nil }
func (method *WrappedBuiltInMethod) Ancestors() []EmeraldValue { return []EmeraldValue{} }
func (method *WrappedBuiltInMethod) HashKey() string           { return method.Inspect() }
