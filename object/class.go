package object

import (
	"reflect"
)

type Class struct {
	*BaseEmeraldValue
	Name        string
	StaticClass *StaticClass
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
	builtInMethodSet,
	staticBuiltInMethodSet BuiltInMethodSet,
	modules ...EmeraldValue,
) *Class {
	var staticParent *StaticClass

	if parentClass != nil {
		staticParent = parentClass.StaticClass
	}

	class := &Class{
		Name:        name,
		parentClass: parentClass,
		BaseEmeraldValue: &BaseEmeraldValue{
			builtInMethodSet: builtInMethodSet,
			includedModules:  modules,
		},
	}

	class.StaticClass = NewStaticClass(name, class, staticBuiltInMethodSet, staticParent)

	for _, module := range modules {
		class.StaticClass.Include(module.(*Module).StaticModule)
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
