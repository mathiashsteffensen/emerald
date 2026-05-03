package core

import (
	"emerald/object"
	"strings"
)

func (rt *Runtime) InitClass() {
	rt.Class = object.NewClass("Class", nil, nil, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	rt.DefineSingletonMethod(rt.Class, "new", rt.classSingletonNew())

	rt.DefineMethod(rt.Class, "new", rt.classNew())
	rt.DefineMethod(rt.Class, "name", rt.className())
	rt.DefineMethod(rt.Class, "ancestors", rt.classAncestors())
}

func (rt *Runtime) classAncestors() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewArray(ctx.Self.Ancestors())
	}
}

func (rt *Runtime) className() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		var namespaces strings.Builder

		parent := ctx.Self.ParentNamespace()
		for parent != nil &&
			parent != rt.Object &&
			(parent.Type() == object.CLASS_VALUE || parent.Type() == object.MODULE_VALUE) {

			switch parent := parent.(type) {
			case *object.Module:
				namespaces.WriteString(parent.Name)
			case *object.Class:
				namespaces.WriteString(parent.Name)
			}

			namespaces.WriteString("::")

			parent = parent.ParentNamespace()
		}

		switch self := ctx.Self.(type) {
		case *object.Module:
			namespaces.WriteString(self.Name)
		case *object.Class:
			namespaces.WriteString(self.Name)
		}

		return rt.NewString(namespaces.String())
	}
}

func (rt *Runtime) classNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		instance := ctx.Self.(*object.Class).New()

		if instance.RespondsTo("initialize", instance) {
			rt.Send(instance, "initialize", ctx.Block, kwargs, args...)
		}

		return instance
	}
}

func (rt *Runtime) classSingletonNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return object.NewClass("", rt.Object, rt.Class, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
	}
}
