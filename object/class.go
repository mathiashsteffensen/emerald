package object

import (
	"reflect"
)

type Class struct {
	*BaseEmeraldValue
	Name                   string
	baseClass              EmeraldValue
	singleton              *SingletonClass
	super                  EmeraldValue
	staticBuiltInMethodSet BuiltInMethodSet
}

func (c *Class) Type() EmeraldValueType { return CLASS_VALUE }
func (c *Class) Inspect() string {
	return c.Name
}
func (c *Class) Class() EmeraldValue {
	if c.singleton != nil {
		return c.singleton
	}
	return c.baseClass
}
func (c *Class) SingletonClass() EmeraldValue {
	if c.singleton == nil {
		c.singleton = NewSingletonClass(c, c.baseClass, &BaseEmeraldValue{builtInMethodSet: c.staticBuiltInMethodSet})
	}
	return c.singleton
}
func (c *Class) Super() EmeraldValue       { return c.super }
func (c *Class) SetSuper(val EmeraldValue) { c.super = val }
func (c *Class) Ancestors() []EmeraldValue {
	ancestors := []EmeraldValue{c}
	ancestors = append(ancestors, c.IncludedModules()...)

	super := c.Super()

	if super != nil && !reflect.ValueOf(super).IsNil() {
		ancestors = append(ancestors, super.(*Class).Ancestors()...)
	}

	return ancestors
}
func (c *Class) HashKey() string { return c.Inspect() }

func (c *Class) New() *Instance {
	return &Instance{
		BaseEmeraldValue: &BaseEmeraldValue{},
		baseClass:        c,
	}
}

func NewClass(
	name string,
	super *Class,
	staticParent EmeraldValue,
	builtInMethodSet,
	staticBuiltInMethodSet BuiltInMethodSet,
	modules ...EmeraldValue,
) *Class {
	class := &Class{
		Name:  name,
		super: super,
		BaseEmeraldValue: &BaseEmeraldValue{
			builtInMethodSet: builtInMethodSet,
			includedModules:  modules,
		},
		baseClass:              staticParent,
		staticBuiltInMethodSet: staticBuiltInMethodSet,
	}

	if len(staticBuiltInMethodSet) != 0 {
		class.SingletonClass()
	}

	return class
}
