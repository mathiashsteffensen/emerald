package object

type Module struct {
	*BaseEmeraldValue
	Name  string
	class EmeraldValue
}

var Modules = map[string]*Module{}

func (m *Module) Type() EmeraldValueType { return MODULE_VALUE }
func (m *Module) Inspect() string {
	return m.Name
}
func (m *Module) Class() EmeraldValue { return m.class }
func (m *Module) Super() EmeraldValue { return nil }
func (m *Module) Ancestors() []EmeraldValue {
	ancestors := []EmeraldValue{m}
	ancestors = append(ancestors, m.IncludedModules()...)

	return ancestors
}
func (m *Module) HashKey() string { return m.Inspect() }

func NewModule(name string, builtInMethodSet, staticBuiltInMethodSet BuiltInMethodSet, parentClass EmeraldValue, modules ...EmeraldValue) *Module {
	mod := &Module{
		BaseEmeraldValue: &BaseEmeraldValue{
			builtInMethodSet: builtInMethodSet,
			includedModules:  modules,
		},
		Name: name,
	}

	mod.class = NewSingletonClass(mod, staticBuiltInMethodSet, parentClass)

	if name != "" {
		Modules[name] = mod
	}

	return mod
}
