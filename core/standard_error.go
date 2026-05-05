package core

import (
	"emerald/object"
	"fmt"
)

func (rt *Runtime) InitStandardError() {
	rt.StandardError = rt.DefineClass("StandardError", rt.Exception)

	rt.DefineSingletonMethod(rt.StandardError, "new", rt.exceptionNew(rt.NewStandardError))
}

type StandardErrorInstance struct {
	*object.Instance
	message string
}

func (err *StandardErrorInstance) Inspect() string {
	return fmt.Sprintf("#<StandardError: %s>", err.message)
}

func (err *StandardErrorInstance) Message() string {
	return err.message
}

func (err *StandardErrorInstance) ClassName() string {
	return "StandardError"
}

func (rt *Runtime) NewStandardError(msg string) object.EmeraldError {
	return &StandardErrorInstance{rt.StandardError.Heap.(*object.Class).New(), msg}
}
