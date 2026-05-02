package core

import (
	"emerald/object"
	"fmt"
)

func (rt *Runtime) InitNameError() {
	rt.NameError = rt.DefineClass("NameError", rt.StandardError)

	rt.DefineSingletonMethod(rt.NameError, "new", rt.exceptionNew(rt.NewNameError))
}

type NameErrorInstance struct {
	*object.Instance
	message string
}

func (err *NameErrorInstance) Inspect() string {
	return fmt.Sprintf("#<NameError: %s>", err.message)
}

func (err *NameErrorInstance) Message() string {
	return err.message
}

func (err *NameErrorInstance) ClassName() string {
	return "NameError"
}

func (rt *Runtime) NewNameError(msg string) object.EmeraldError {
	return &NameErrorInstance{rt.NameError.New(), msg}
}
