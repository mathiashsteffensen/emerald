package object

import "fmt"

type ReturnValue struct {
	*BaseEmeraldValue
	Value EmeraldValue
}

func (rv *ReturnValue) Class() EmeraldValue       { return EmeraldValue{} }
func (rv *ReturnValue) Super() EmeraldValue       { return EmeraldValue{} }
func (rv *ReturnValue) Ancestors() []EmeraldValue { return []EmeraldValue{} }
func (rv *ReturnValue) Type() EmeraldValueType    { return RETURN_VALUE }
func (rv *ReturnValue) Inspect() string {
	return fmt.Sprintf("return %s", rv.Value.Inspect())
}
func (rv *ReturnValue) HashKey() string              { return rv.Inspect() }
func (rv *ReturnValue) SingletonClass() EmeraldValue { return EmeraldValue{} }
