package object

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
func (c *Class) New() *Instance {
	return &Instance{class: c, BaseEmeraldValue: c.BaseEmeraldValue, BuiltInSingletonMethods: BuiltInMethodSet{}}
}

func NewClass(
	name string,
	parentClass *Class,
	builtInMethodSet BuiltInMethodSet,
	staticBuiltInMethodSet BuiltInMethodSet,
) *Class {
	class := &Class{
		Name:        name,
		parentClass: parentClass,
		BaseEmeraldValue: &BaseEmeraldValue{
			builtInMethodSet:       builtInMethodSet,
			staticBuiltInMethodSet: staticBuiltInMethodSet,
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
