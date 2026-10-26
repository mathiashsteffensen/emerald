package object

import (
	"fmt"
	"strings"
)

type Instance struct {
	*BaseEmeraldValue
	baseClass EmeraldValue
	singleton *SingletonClass
	backtrace []string
}

func (i *Instance) Type() EmeraldValueType { return INSTANCE_VALUE }
func (i *Instance) Class() EmeraldValue {
	if i.singleton != nil {
		return NewHeapObject(i.singleton)
	}
	return i.baseClass
}
func (i *Instance) SingletonClass() EmeraldValue {
	if i.singleton == nil {
		i.singleton = NewSingletonClass(NewHeapObject(i), i.baseClass, i.BaseEmeraldValue)
	}
	return NewHeapObject(i.singleton)
}
func (i *Instance) Super() EmeraldValue { return EmeraldValue{} }
func (i *Instance) Ancestors() []EmeraldValue {
	return append(i.Class().Ancestors(), NewHeapObject(i))
}
func (i *Instance) Include(mod EmeraldValue) {
	i.SingletonClass().Include(mod)
}
func (i *Instance) DefineMethod(block EmeraldValue, args ...EmeraldValue) {
	i.SingletonClass().DefineMethod(block, args...)
}
func (i *Instance) Inspect() string {
	var out strings.Builder

	out.WriteString("#<")

	realClass := RealClass(NewHeapObject(i))
	if class, ok := realClass.Heap.(*Class); ok {
		out.WriteString(class.Name)
	} else {
		out.WriteString("Unknown")
	}

	out.WriteString(":")
	out.WriteString(fmt.Sprintf("%p", i))
	out.WriteString(">")

	return out.String()
}
func (i *Instance) HashKey() string { return i.Inspect() }

// Backtrace returns the call stack captured when this instance was raised as an error, or nil
func (i *Instance) Backtrace() []string {
	if i == nil {
		return nil
	}
	return i.backtrace
}

func (i *Instance) SetBacktrace(backtrace []string) {
	if i == nil {
		return
	}
	i.backtrace = backtrace
}
