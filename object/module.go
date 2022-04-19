package object

import "reflect"

type Module struct {
	*BaseEmeraldValue
	Name string
}

var Modules = map[string]*Module{}

func (c *Module) Type() EmeraldValueType { return MODULE_VALUE }
func (c *Module) Inspect() string {
	return c.Name
}
func (c *Module) ParentClass() EmeraldValue { return nil }
func (c *Module) Ancestors() []EmeraldValue {
	ancestors := []EmeraldValue{c}
	ancestors = append(ancestors, c.IncludedModules()...)

	super := c.ParentClass()
	reflected := reflect.ValueOf(super)
	if super != nil && reflected.IsValid() && !reflected.IsNil() {
		ancestors = append(ancestors, super.Ancestors()...)
	}

	return ancestors
}

func NewModule(name string, builtInMethodSet, staticBuiltInMethodSet BuiltInMethodSet, modules ...EmeraldValue) *Module {
	mod := &Module{
		BaseEmeraldValue: &BaseEmeraldValue{
			staticBuiltInMethodSet: staticBuiltInMethodSet,
			builtInMethodSet:       builtInMethodSet,
			includedModules:        modules,
		},
		Name: name,
	}

	if name != "" {
		Modules[name] = mod
	}

	return mod
}
