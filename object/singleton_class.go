package object

import (
	"fmt"
	"reflect"
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
	ancestors := []EmeraldValue{c}
	ancestors = append(ancestors, c.IncludedModules()...)

	super := c.Super()
	reflected := reflect.ValueOf(super)
	if super != nil && reflected.IsValid() && !reflected.IsNil() {
		ancestors = append(ancestors, super.Ancestors()...)
	}

	return ancestors
}
func (c *SingletonClass) HashKey() string { return c.Inspect() }

func NewSingletonClass(instance EmeraldValue, set BuiltInMethodSet, super EmeraldValue) *SingletonClass {
	return &SingletonClass{
		BaseEmeraldValue: &BaseEmeraldValue{builtInMethodSet: set},
		Instance:         instance,
		super:            super,
	}
}
