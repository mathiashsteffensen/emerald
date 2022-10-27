package core

import (
	"strings"
	"emerald/object"
)

var Class *object.Class

func InitClass() {
	Class = object.NewClass("Class", nil, nil, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	DefineSingletonMethod(Class, "new", classSingletonNew())

	DefineMethod(Class, "new", classNew())
	DefineMethod(Class, "name", className())
	DefineMethod(Class, "ancestors", classAncestors())
}

func classAncestors() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewArray(ctx.Self.Ancestors())
	}
}

func className() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		var namespaces strings.Builder

		parent := ctx.Self.ParentNamespace()
		for parent != nil &&
			parent != Object &&
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

		namespaces.WriteString(ctx.Self.(*object.Class).Name)

		return NewString(namespaces.String())
	}
}

func classNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		instance := ctx.Self.(*object.Class).New()

		if instance.RespondsTo("initialize", instance) {
			Send(instance, "initialize", ctx.Block, args...)
		}

		return instance
	}
}

func classSingletonNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return object.NewClass("", Object, Object.Class(), object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
	}
}
