package object

import (
	"bytes"
	"fmt"
)

type Instance struct {
	*BaseEmeraldValue
	class EmeraldValue
}

func (i *Instance) Type() EmeraldValueType    { return INSTANCE_VALUE }
func (i *Instance) Class() EmeraldValue       { return i.class }
func (i *Instance) Super() EmeraldValue       { return nil }
func (i *Instance) Ancestors() []EmeraldValue { return append(i.class.Ancestors(), i) }
func (i *Instance) Inspect() string {
	var out bytes.Buffer

	out.WriteString("#<")
	out.WriteString(i.Class().Super().(*Class).Name)
	out.WriteString(":")
	out.WriteString(fmt.Sprintf("%p", i))
	out.WriteString(">")

	return out.String()
}
func (i *Instance) HashKey() string { return i.Inspect() }
