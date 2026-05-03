package object

import (
	"fmt"
	"strings"
)

type Instance struct {
	*BaseEmeraldValue
	baseClass EmeraldValue
	singleton *SingletonClass
}

func (i *Instance) Type() EmeraldValueType { return INSTANCE_VALUE }
func (i *Instance) Class() EmeraldValue {
	if i.singleton != nil {
		return i.singleton
	}
	return i.baseClass
}
func (i *Instance) SingletonClass() EmeraldValue {
	if i.singleton == nil {
		i.singleton = NewSingletonClass(i, i.baseClass, i.BaseEmeraldValue)
	}
	return i.singleton
}
func (i *Instance) Super() EmeraldValue { return nil }
func (i *Instance) Ancestors() []EmeraldValue {
	return append(i.Class().Ancestors(), i)
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

	out.WriteString(RealClass(i).(*Class).Name)

	out.WriteString(":")
	out.WriteString(fmt.Sprintf("%p", i))
	out.WriteString(">")

	return out.String()
}
func (i *Instance) HashKey() string { return i.Inspect() }
