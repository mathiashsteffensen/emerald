package core

import (
	"bytes"
	"emerald/object"
)

var Class *object.Class

func InitClass() {
	Class = object.NewClass(
		"Class",
		nil,
		nil,
		object.BuiltInMethodSet{
			"ancestors": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewArray(target.Ancestors())
			},
			"methods": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				methods := []object.EmeraldValue{}

				for _, method := range target.Methods(target) {
					methods = append(methods, NewSymbol(method))
				}

				return NewArray(methods)
			},
			"name": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				self := target.(*object.Class)

				var namespaces bytes.Buffer

				parent := self.ParentNamespace()
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

				namespaces.WriteString(self.Name)

				return NewString(namespaces.String())
			},
			"new": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, _yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return target.(*object.Class).New()
			},
		},
		object.BuiltInMethodSet{
			"new": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return object.NewClass("", Object, Object.Class(), object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
			},
		},
	)
}
