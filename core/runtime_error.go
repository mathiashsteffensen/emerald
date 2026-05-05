package core

import (
	"emerald/object"
	"fmt"
)

func (rt *Runtime) InitRuntimeError() {
	rt.RuntimeError = rt.DefineClass("RuntimeError", rt.StandardError)

	rt.DefineSingletonMethod(rt.RuntimeError, "new", rt.exceptionNew(rt.NewRuntimeError))
}

type RuntimeErrorInstance struct {
	*object.Instance
	message string
}

func (err *RuntimeErrorInstance) Inspect() string {
	return fmt.Sprintf("#<RuntimeError: %s>", err.message)
}

func (err *RuntimeErrorInstance) Message() string {
	return err.message
}

func (err *RuntimeErrorInstance) ClassName() string {
	return "RuntimeError"
}

func (rt *Runtime) NewRuntimeError(msg string) object.EmeraldError {
	return &RuntimeErrorInstance{rt.RuntimeError.Heap.(*object.Class).New(), msg}
}
