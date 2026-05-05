package object

import (
	"fmt"
)

type SingletonClass struct {
	*BaseEmeraldValue
	Instance EmeraldValue
	super    EmeraldValue
}

func (c *SingletonClass) Type() EmeraldValueType    { return STATIC_CLASS_VALUE }
func (c *SingletonClass) Inspect() string           { return fmt.Sprintf("#<Class:%s>", c.Instance.Inspect()) }
func (c *SingletonClass) Class() EmeraldValue       { return c.super }
func (c *SingletonClass) Super() EmeraldValue       { return c.super }
func (c *SingletonClass) SetSuper(val EmeraldValue) { c.super = val }
func (c *SingletonClass) Ancestors() []EmeraldValue {
	ancestors := []EmeraldValue{NewHeapObject(c)}
	ancestors = append(ancestors, c.IncludedModules()...)

	super := c.Super()
	if !super.IsNil() {
		ancestors = append(ancestors, super.Ancestors()...)
	}

	return ancestors
}
func (c *SingletonClass) HashKey() string              { return c.Inspect() }
func (c *SingletonClass) SingletonClass() EmeraldValue { return EmeraldValue{} }

func NewSingletonClass(instance EmeraldValue, super EmeraldValue, base *BaseEmeraldValue) *SingletonClass {
	return &SingletonClass{
		BaseEmeraldValue: base,
		Instance:         instance,
		super:            super,
	}
}
