package core

import (
	"emerald/object"

	"github.com/dlclark/regexp2"
)

// Inherited native new keeps Class#new's subclass initialization policy, not
// the owner's argument-specific constructor, while allocating native storage.
func (rt *Runtime) defineNativeConstructor(owner object.EmeraldValue, constructor object.BuiltInMethod) {
	rt.DefineSingletonMethod(owner, "new", func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if ctx.Self == owner {
			return constructor(ctx, kwargs, args...)
		}

		base := ctx.Self.Heap.(*object.Class).New()
		var native object.HeapObject = base
		switch owner {
		case rt.String:
			native = &StringInstance{Instance: base}
		case rt.Regexp:
			native = &RegexpInstance{Instance: base, Expression: regexp2.MustCompile("", 0)}
		case rt.Range:
			native = &RangeInstance{Instance: base, Begin: rt.NULL, End: rt.NULL}
		case rt.Time:
			native = &TimeInstance{Instance: base}
		case rt.IO:
			native = &IOInstance{Instance: base, FileDescriptor: ^uintptr(0), Closed: true}
		case rt.TCPServer:
			native = &TCPServerInstance{Instance: base}
		case rt.Exception, rt.StandardError, rt.ArgumentError, rt.TypeError, rt.NameError, rt.NoMethodError, rt.RuntimeError, rt.LoadError:
			native = &ExceptionInstance{Instance: base}
		case rt.Class:
			native = object.NewClass("", rt.Object.Heap.(*object.Class), ctx.Self, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
			native.SingletonClass().Heap.(*object.SingletonClass).SetSuper(ctx.Self)
		}
		instance := object.NewHeapObject(native)
		if instance.RespondsTo("initialize", instance) {
			rt.Send(instance, "initialize", ctx.Block, kwargs, args...)
		}
		return instance
	})
}
