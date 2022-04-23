package object

import "reflect"

type StaticClass struct {
	*BaseEmeraldValue
	Name        string
	Class       EmeraldValue
	parentClass EmeraldValue
}

func (c *StaticClass) Type() EmeraldValueType { return STATIC_CLASS_VALUE }
func (c *StaticClass) Inspect() string {
	return c.Name
}
func (c *StaticClass) ParentClass() EmeraldValue { return c.parentClass }
func (c *StaticClass) Ancestors() []EmeraldValue {
	ancestors := []EmeraldValue{c}
	ancestors = append(ancestors, c.IncludedModules()...)

	super := c.ParentClass()
	reflected := reflect.ValueOf(super)
	if super != nil && reflected.IsValid() && !reflected.IsNil() {
		ancestors = append(ancestors, super.Ancestors()...)
	}

	return ancestors
}

func NewStaticClass(name string, class EmeraldValue, set BuiltInMethodSet, parentClass *StaticClass) *StaticClass {
	return &StaticClass{
		BaseEmeraldValue: &BaseEmeraldValue{builtInMethodSet: set},
		Name:             name,
		Class:            class,
		parentClass:      parentClass,
	}
}
