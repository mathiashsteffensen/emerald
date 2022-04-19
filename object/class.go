package object

import "reflect"

type Class struct {
	*BaseEmeraldValue
	Name        string
	parentClass EmeraldValue
}

var Classes = map[string]*Class{}

func (c *Class) Type() EmeraldValueType { return CLASS_VALUE }
func (c *Class) Inspect() string {
	return c.Name
}
func (c *Class) ParentClass() EmeraldValue { return c.parentClass }
func (c *Class) Ancestors() []EmeraldValue {
	ancestors := []EmeraldValue{c}
	ancestors = append(ancestors, c.IncludedModules()...)

	super := c.ParentClass()
	reflected := reflect.ValueOf(super)
	if super != nil && reflected.IsValid() && !reflected.IsNil() {
		ancestors = append(ancestors, super.Ancestors()...)
	}

	return ancestors
}

func (c *Class) New() *Instance {
	return &Instance{class: c, BaseEmeraldValue: c.BaseEmeraldValue, BuiltInSingletonMethods: BuiltInMethodSet{}}
}

func NewClass(
	name string,
	parentClass *Class,
	builtInMethodSet BuiltInMethodSet,
	staticBuiltInMethodSet BuiltInMethodSet,
	modules ...EmeraldValue,
) *Class {
	class := &Class{
		Name:        name,
		parentClass: parentClass,
		BaseEmeraldValue: &BaseEmeraldValue{
			builtInMethodSet:       builtInMethodSet,
			staticBuiltInMethodSet: staticBuiltInMethodSet,
			includedModules:        modules,
		},
	}

	if name != "" {
		Classes[name] = class
	}

	return class
}

func GetClassByName(name string) (*Class, bool) {
	class, ok := Classes[name]

	return class, ok
}
