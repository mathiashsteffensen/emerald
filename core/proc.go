package core

import "emerald/object"

// ProcInstance is a block captured as a first-class value, as produced by a
// method's &block parameter or Proc.new.
type ProcInstance struct {
	*object.Instance
	Block object.EmeraldValue
}

func (rt *Runtime) InitProc() {
	rt.Proc = rt.DefineClass("Proc", rt.Object)

	// Proc.new takes the block of the call, like Ruby, instead of inheriting
	// Class#new, which knows nothing about blocks.
	rt.DefineSingletonMethod(rt.Proc, "new", rt.procNew())
	rt.DefineMethod(rt.Proc, "initialize", rt.procInitialize())
	rt.DefineMethod(rt.Proc, "call", rt.procCall())
}

func (rt *Runtime) NewProc(block object.EmeraldValue) object.EmeraldValue {
	return object.NewHeapObject(&ProcInstance{
		Instance: rt.Proc.Heap.(*object.Class).New(),
		Block:    block,
	})
}

// procNew wraps the block passed to Proc.new in an instance of the receiver
// class. A subclass gets its own class and, like Class#new, its initialize
// receives the call's block and arguments. Direct Proc.new has Proc's
// zero-argument initialize.
func (rt *Runtime) procNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if ctx.Block.IsNil() {
			err := rt.newArgumentError("tried to create Proc object without a block")
			rt.Raise(err)
			return object.NewHeapObject(err)
		}

		instance := object.NewHeapObject(&ProcInstance{
			Instance: ctx.Self.Heap.(*object.Class).New(),
			Block:    ctx.Block,
		})

		if instance.RespondsTo("initialize", instance) {
			rt.Send(instance, "initialize", ctx.Block, kwargs, args...)
		}

		return instance
	}
}

// procInitialize keeps Proc.new's zero-argument contract. A subclass that
// defines initialize replaces it and receives the original arguments.
func (rt *Runtime) procInitialize() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 0, 0); err != nil {
			return object.NewHeapObject(err)
		}

		return rt.NULL
	}
}

// procCall runs the block with the arguments the caller supplies. A block
// keeps its defining context, so captured locals and self stay available.
func (rt *Runtime) procCall() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		proc, ok := ctx.Self.Heap.(*ProcInstance)
		if !ok || proc.Block.IsNil() {
			return object.NewHeapObject(rt.Raise(rt.NewTypeError("Proc#call requires a block")))
		}

		return rt.EvalBlock(proc.Block, kwargs, args...)
	}
}
