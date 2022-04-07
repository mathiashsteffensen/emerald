package object

type Class struct {
	*BaseEmeraldValue
	Name        string
	parentClass EmeraldValue
}

func (c *Class) Type() EmeraldValueType { return CLASS_VALUE }
func (c *Class) Inspect() string {
	return c.Name
}
func (c *Class) ParentClass() EmeraldValue { return c.parentClass }
func (c *Class) New() *Instance {
	return &Instance{class: c, BaseEmeraldValue: c.BaseEmeraldValue}
}

func NewClass(name string, parentClass *Class, builtInMethodSet BuiltInMethodSet, staticBuiltInMethodSet BuiltInMethodSet) *Class {
	class := &Class{
		Name:        name,
		parentClass: parentClass,
		BaseEmeraldValue: &BaseEmeraldValue{
			builtInMethodSet:       builtInMethodSet,
			staticBuiltInMethodSet: staticBuiltInMethodSet,
		},
	}

	defaultEnvironment.Set(name, class)

	return class
}
