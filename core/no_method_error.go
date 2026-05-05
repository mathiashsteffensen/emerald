package core

import "emerald/object"

type NoMethodErrorInstance struct {
	*object.Instance
	message string
}

func (err *NoMethodErrorInstance) Message() string   { return err.message }
func (err *NoMethodErrorInstance) ClassName() string { return "NoMethodError" }

func (rt *Runtime) InitNoMethodError() {
	rt.NoMethodError = rt.DefineClass("NoMethodError", rt.StandardError)

	rt.DefineSingletonMethod(rt.NoMethodError, "new", rt.exceptionNew(rt.NewNoMethodError))
}

func (rt *Runtime) NewNoMethodError(msg string) object.EmeraldError {
	return &NoMethodErrorInstance{
		Instance: rt.NoMethodError.Heap.(*object.Class).New(),
		message:  msg,
	}
}
