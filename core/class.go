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
			"ancestors": func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewArray(self.Ancestors())
			},
			"methods": func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				methods := []object.EmeraldValue{}

				for _, method := range self.Methods(self) {
					methods = append(methods, NewSymbol(method))
				}

				return NewArray(methods)
			},
			"name": func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
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

				namespaces.WriteString(self.(*object.Class).Name)

				return NewString(namespaces.String())
			},
			"new": func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				instance := self.(*object.Class).New()

				if instance.RespondsTo("initialize", instance) {
					Send(instance, "initialize", block, args...)
				}

				return instance
			},
		},
		object.BuiltInMethodSet{
			"new": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return object.NewClass("", Object, Object.Class(), object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
			},
		},
	)
}
