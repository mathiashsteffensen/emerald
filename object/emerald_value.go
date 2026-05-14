package object

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type MethodVisibility string

const (
	PUBLIC    MethodVisibility = "public"
	PRIVATE   MethodVisibility = "private"
	PROTECTED MethodVisibility = "protected"
)

type (
	// BuiltInMethod - The type signature of an Emerald method defined in Go compiler
	BuiltInMethod func(ctx *Context, kwargs map[string]EmeraldValue, args ...EmeraldValue) EmeraldValue

	// WrappedBuiltInMethod -  Wraps a built-in method so that it conforms to the EmeraldValue interface
	WrappedBuiltInMethod struct {
		*BaseEmeraldValue
		Method     BuiltInMethod
		Visibility MethodVisibility
	}

	// BuiltInMethodSet - Stores an objects built-in method set
	BuiltInMethodSet map[string]*WrappedBuiltInMethod

	// DefinedMethodSet - Stores an objects methods defined by the program
	DefinedMethodSet map[string]*ClosedBlock

	EmeraldValueType uint8

	// EmeraldValue **everything** in Emerald is an EmeraldValue
	EmeraldValue struct {
		TypeID EmeraldValueType
		Heap   HeapObject
		Num    uint64
	}

	// HeapObject - All Emerald objects must implement this interface
	HeapObject interface {
		SetType(t EmeraldValueType)
		SetHeap(h HeapObject)
		Type() EmeraldValueType
		Inspect() string
		Class() EmeraldValue
		Super() EmeraldValue
		Ancestors() []EmeraldValue
		IncludedModules() []EmeraldValue
		Include(mod EmeraldValue)
		BuiltInMethodSet() BuiltInMethodSet
		DefinedMethodSet() DefinedMethodSet
		ExtractMethod(name string, extractFrom EmeraldValue, self EmeraldValue) (
			EmeraldValue, // The actual method
			MethodVisibility,
			bool, // Boolean that is true if this method is defined directly on self
			error, // error if no method was found
		)
		Methods() []string
		InstanceVariableGet(name string, extractFrom EmeraldValue, self EmeraldValue) EmeraldValue
		InstanceVariableSet(name string, value EmeraldValue)
		ParentNamespace() EmeraldValue
		SetParentNamespace(parent EmeraldValue)
		NamespaceDefinitionGet(name string) EmeraldValue
		NamespaceDefinitionSet(name string, value EmeraldValue)
		NamespaceDefinitions() map[string]EmeraldValue
		HashKey() string
		SingletonClass() EmeraldValue
		DefineMethod(block EmeraldValue, args ...EmeraldValue)
	}
)

func (v EmeraldValue) IsImmediate() bool {
	switch v.TypeID {
	case INTEGER_VALUE, FLOAT_VALUE, NIL_VALUE, TRUE_VALUE, FALSE_VALUE:
		return true
	default:
		return false
	}
}

func NewHeapObject(obj HeapObject) EmeraldValue {
	reflected := reflect.ValueOf(obj)
	if obj == nil || (reflected.Kind() == reflect.Ptr && reflected.IsNil()) {
		return EmeraldValue{TypeID: NIL_VALUE}
	}
	return EmeraldValue{TypeID: obj.Type(), Heap: obj}
}

func (v EmeraldValue) Type() EmeraldValueType { return v.TypeID }

func (v EmeraldValue) IsNil() bool { return v.Is(NIL_VALUE) }

func (v EmeraldValue) IsFalse() bool { return v.Is(FALSE_VALUE) }

func (v EmeraldValue) IsDefined() bool {
	return v.TypeID != NIL_VALUE || v.Heap != nil
}

func (v EmeraldValue) Is(t EmeraldValueType) bool { return v.TypeID == t }

func (v EmeraldValue) Inspect() string {
	switch v.TypeID {
	case NIL_VALUE:
		return "nil"
	case TRUE_VALUE:
		return "true"
	case FALSE_VALUE:
		return "false"
	case INTEGER_VALUE:
		return strconv.FormatInt(int64(v.Num), 10)
	case FLOAT_VALUE:
		return strconv.FormatFloat(math.Float64frombits(v.Num), 'g', -1, 64)
	default:
		if v.Heap != nil {
			return v.Heap.Inspect()
		}
		return "nil"
	}
}

func (v EmeraldValue) Class() EmeraldValue {
	if v.IsImmediate() {
		return NewHeapObject(v.Heap)
	}
	if v.Heap != nil {
		return v.Heap.Class()
	}
	return EmeraldValue{}
}

func (v EmeraldValue) Super() EmeraldValue {
	if v.Heap != nil {
		return v.Heap.Super()
	}
	return EmeraldValue{}
}

func (v EmeraldValue) Ancestors() []EmeraldValue {
	if v.IsImmediate() {
		if v.Heap != nil {
			return append(v.Heap.Ancestors(), v)
		}
		return []EmeraldValue{v}
	}
	if v.Heap != nil {
		return v.Heap.Ancestors()
	}
	return []EmeraldValue{}
}

func (v EmeraldValue) IncludedModules() []EmeraldValue {
	if v.Heap != nil {
		return v.Heap.IncludedModules()
	}
	return []EmeraldValue{}
}

func (v EmeraldValue) Include(mod EmeraldValue) {
	if v.Heap != nil {
		v.Heap.Include(mod)
	}
}

func (v EmeraldValue) BuiltInMethodSet() BuiltInMethodSet {
	if v.Heap != nil {
		return v.Heap.BuiltInMethodSet()
	}
	return BuiltInMethodSet{}
}

func (v EmeraldValue) DefinedMethodSet() DefinedMethodSet {
	if v.Heap != nil {
		return v.Heap.DefinedMethodSet()
	}
	return DefinedMethodSet{}
}

func (v EmeraldValue) ExtractMethod(name string, extractFrom EmeraldValue, self EmeraldValue) (EmeraldValue, MethodVisibility, bool, error) {
	if v.IsImmediate() {
		if v.Heap != nil {
			return v.Heap.ExtractMethod(name, extractFrom, self)
		}
	}
	if v.Heap != nil {
		return v.Heap.ExtractMethod(name, extractFrom, self)
	}
	return EmeraldValue{}, PUBLIC, false, fmt.Errorf("undefined method %s", name)
}

func (v EmeraldValue) Methods() []string {
	if v.Heap != nil {
		return v.Heap.Methods()
	}
	return []string{}
}

func (v EmeraldValue) RespondsTo(name string, self EmeraldValue) bool {
	if v.Heap != nil {
		_, visibility, _, err := v.Heap.ExtractMethod(name, v.Class(), self)
		return err == nil && visibility == PUBLIC
	}
	return false
}

func (v EmeraldValue) InstanceVariableGet(name string, extractFrom EmeraldValue, self EmeraldValue) EmeraldValue {
	if v.Heap != nil {
		return v.Heap.InstanceVariableGet(name, extractFrom, self)
	}
	return EmeraldValue{}
}

func (v EmeraldValue) InstanceVariableSet(name string, value EmeraldValue) {
	if v.Heap != nil {
		v.Heap.InstanceVariableSet(name, value)
	}
}

func (v EmeraldValue) ParentNamespace() EmeraldValue {
	if v.Heap != nil {
		return v.Heap.ParentNamespace()
	}
	return EmeraldValue{}
}

func (v EmeraldValue) SetParentNamespace(parent EmeraldValue) {
	if v.Heap != nil {
		v.Heap.SetParentNamespace(parent)
	}
}

func (v EmeraldValue) NamespaceDefinitionGet(name string) EmeraldValue {
	if v.Heap != nil {
		return v.Heap.NamespaceDefinitionGet(name)
	}
	return EmeraldValue{}
}

func (v EmeraldValue) NamespaceDefinitionSet(name string, value EmeraldValue) {
	if v.Heap != nil {
		v.Heap.NamespaceDefinitionSet(name, value)
	}
}

func (v EmeraldValue) NamespaceDefinitions() map[string]EmeraldValue {
	if v.Heap != nil {
		return v.Heap.NamespaceDefinitions()
	}
	return map[string]EmeraldValue{}
}

func (v EmeraldValue) HashKey() string {
	switch v.TypeID {
	case NIL_VALUE:
		return "nil"
	case TRUE_VALUE:
		return "true"
	case FALSE_VALUE:
		return "false"
	case INTEGER_VALUE:
		return "int:" + strconv.FormatInt(int64(v.Num), 10)
	case FLOAT_VALUE:
		return "float:" + strconv.FormatFloat(math.Float64frombits(v.Num), 'g', -1, 64)
	default:
		if v.Heap != nil {
			return v.Heap.HashKey()
		}
		return "nil"
	}
}

func (v EmeraldValue) SingletonClass() EmeraldValue {
	if v.Heap != nil {
		return v.Heap.SingletonClass()
	}
	return EmeraldValue{}
}

func (v EmeraldValue) DefineMethod(block EmeraldValue, args ...EmeraldValue) {
	if v.Heap != nil {
		v.Heap.DefineMethod(block, args...)
	}
}

const (
	NIL_VALUE EmeraldValueType = iota
	CLASS_VALUE
	STATIC_CLASS_VALUE
	MODULE_VALUE
	INSTANCE_VALUE
	BLOCK_VALUE
	RETURN_VALUE
	INTEGER_VALUE
	FLOAT_VALUE
	TRUE_VALUE
	FALSE_VALUE
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
	case INTEGER_VALUE:
		return "Integer"
	case FLOAT_VALUE:
		return "Float"
	case TRUE_VALUE:
		return "true"
	case FALSE_VALUE:
		return "false"
	case NIL_VALUE:
		return "nil"
	}

	return ""
}

func (method *WrappedBuiltInMethod) Inspect() string                                       { return fmt.Sprintf("#<Block:%p>", method) }
func (method *WrappedBuiltInMethod) Type() EmeraldValueType                                { return BLOCK_VALUE }
func (method *WrappedBuiltInMethod) Class() EmeraldValue                                   { return EmeraldValue{} }
func (method *WrappedBuiltInMethod) Super() EmeraldValue                                   { return EmeraldValue{} }
func (method *WrappedBuiltInMethod) Ancestors() []EmeraldValue                             { return []EmeraldValue{} }
func (method *WrappedBuiltInMethod) HashKey() string                                       { return method.Inspect() }
func (method *WrappedBuiltInMethod) SingletonClass() EmeraldValue                          { return EmeraldValue{} }
func (method *WrappedBuiltInMethod) SetType(t EmeraldValueType)                            {}
func (method *WrappedBuiltInMethod) SetHeap(h HeapObject)                                  {}
func (method *WrappedBuiltInMethod) Include(mod EmeraldValue)                              {}
func (method *WrappedBuiltInMethod) DefineMethod(block EmeraldValue, args ...EmeraldValue) {}

func RealClass(val EmeraldValue) EmeraldValue {
	class := val.Class()
	for !class.IsNil() && class.Is(STATIC_CLASS_VALUE) {
		class = class.Super()
	}
	return class
}
