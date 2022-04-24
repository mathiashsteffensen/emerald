package object

import (
	"reflect"
)

type Class struct {
	*BaseEmeraldValue
	Name  string
	class EmeraldValue
	super EmeraldValue
}

var Classes = map[string]*Class{}

func (c *Class) Type() EmeraldValueType { return CLASS_VALUE }
func (c *Class) Inspect() string {
	return c.Name
}
func (c *Class) Class() EmeraldValue       { return c.class }
func (c *Class) Super() EmeraldValue       { return c.super }
func (c *Class) SetSuper(val EmeraldValue) { c.super = val }
func (c *Class) Ancestors() []EmeraldValue {
	ancestors := []EmeraldValue{c}
	ancestors = append(ancestors, c.IncludedModules()...)

	super := c.Super()
	reflected := reflect.ValueOf(super)
	if super != nil && reflected.IsValid() && !reflected.IsNil() {
		ancestors = append(ancestors, super.Ancestors()...)
	}

	return ancestors
}

func (c *Class) New() *Instance {
	instance := &Instance{}

	singleton := NewSingletonClass(instance, BuiltInMethodSet{}, c)

	instance.class = singleton
	instance.BaseEmeraldValue = singleton.BaseEmeraldValue

	return instance
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
	}

	class.class = NewSingletonClass(class, staticBuiltInMethodSet, staticParent)

	if name != "" {
		Classes[name] = class
	}

	return class
}

func GetClassByName(name string) (*Class, bool) {
	class, ok := Classes[name]

	return class, ok
}
