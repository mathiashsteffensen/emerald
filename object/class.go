package object

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
	return c.SingletonClass()
}
func (c *Class) SingletonClass() EmeraldValue {
	if c.singleton == nil {
		super := c.baseClass
		if !c.super.IsNil() {
			super = c.super.SingletonClass()
		}
		c.singleton = NewSingletonClass(NewHeapObject(c), super, &BaseEmeraldValue{builtInMethodSet: c.staticBuiltInMethodSet})
	}
	return NewHeapObject(c.singleton)
}
func (c *Class) Super() EmeraldValue       { return c.super }
func (c *Class) SetSuper(val EmeraldValue) { c.super = val }
func (c *Class) Ancestors() []EmeraldValue {
	ancestors := []EmeraldValue{NewHeapObject(c)}
	ancestors = append(ancestors, c.IncludedModules()...)

	super := c.Super()

	if !super.IsNil() {
		ancestors = append(ancestors, super.Heap.(*Class).Ancestors()...)
	}

	return ancestors
}
func (c *Class) HashKey() string { return c.Inspect() }

func (c *Class) New() *Instance {
	return &Instance{
		BaseEmeraldValue: &BaseEmeraldValue{},
		baseClass:        NewHeapObject(c),
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
		super: NewHeapObject(super),
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
