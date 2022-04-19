package object

import (
	"bytes"
	"fmt"
)

type Instance struct {
	*BaseEmeraldValue
	class                   *Class
	BuiltInSingletonMethods BuiltInMethodSet
}

func (i *Instance) Type() EmeraldValueType    { return INSTANCE_VALUE }
func (i *Instance) ParentClass() EmeraldValue { return i.class }
func (i *Instance) Ancestors() []EmeraldValue { return i.class.Ancestors() }
func (i *Instance) Inspect() string {
	var out bytes.Buffer

	out.WriteString("#<")
	out.WriteString(i.class.Name)
	out.WriteString(":")
	out.WriteString(fmt.Sprintf("%p", i))
	out.WriteString(">")

	return out.String()
}
