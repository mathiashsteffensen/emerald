package object

import "fmt"

type ReturnValue struct {
	*BaseEmeraldValue
	Value EmeraldValue
}

func (rv *ReturnValue) ParentClass() EmeraldValue { return nil }
func (rv *ReturnValue) Ancestors() []EmeraldValue { return []EmeraldValue{} }
func (rv *ReturnValue) Type() EmeraldValueType    { return RETURN_VALUE }
func (rv *ReturnValue) Inspect() string {
	return fmt.Sprintf("return %s", rv.Value.Inspect())
}
