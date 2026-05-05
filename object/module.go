package object

type Module struct {
	*BaseEmeraldValue
	Name                   string
	baseClass              EmeraldValue
	singleton              *SingletonClass
	staticBuiltInMethodSet BuiltInMethodSet
}

func (m *Module) Type() EmeraldValueType { return MODULE_VALUE }
func (m *Module) Inspect() string {
	return m.Name
}
func (m *Module) Class() EmeraldValue {
	if m.singleton != nil {
		return NewHeapObject(m.singleton)
	}
	return m.baseClass
}
func (m *Module) SingletonClass() EmeraldValue {
	if m.singleton == nil {
		m.singleton = NewSingletonClass(NewHeapObject(m), m.baseClass, &BaseEmeraldValue{builtInMethodSet: m.staticBuiltInMethodSet})
	}
	return NewHeapObject(m.singleton)
}
func (m *Module) Super() EmeraldValue { return EmeraldValue{} }
func (m *Module) Ancestors() []EmeraldValue {
	ancestors := []EmeraldValue{NewHeapObject(m)}
	ancestors = append(ancestors, m.IncludedModules()...)

	return ancestors
}
func (m *Module) HashKey() string { return m.Inspect() }

func NewModule(name string, moduleClass EmeraldValue, builtInMethodSet, staticBuiltInMethodSet BuiltInMethodSet, modules ...EmeraldValue) *Module {
	mod := &Module{
		BaseEmeraldValue: &BaseEmeraldValue{
			builtInMethodSet: builtInMethodSet,
			includedModules:  modules,
		},
		Name:                   name,
		baseClass:              moduleClass,
		staticBuiltInMethodSet: staticBuiltInMethodSet,
	}

	if len(staticBuiltInMethodSet) != 0 {
		mod.SingletonClass()
	}

	return mod
}
