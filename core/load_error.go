package core

import (
	"emerald/object"
	"fmt"
)

func (rt *Runtime) InitLoadError() {
	rt.LoadError = rt.DefineClass("LoadError", rt.StandardError)

	rt.DefineSingletonMethod(rt.LoadError, "new", rt.exceptionNew(rt.NewLoadError))
}

type LoadErrorInstance struct {
	*object.Instance
	message string
}

func (err *LoadErrorInstance) Inspect() string {
	return fmt.Sprintf("#<LoadError: %s>", err.message)
}

func (err *LoadErrorInstance) Message() string {
	return err.message
}

func (err *LoadErrorInstance) ClassName() string {
	return "LoadError"
}

func (rt *Runtime) NewLoadError(msg string) object.EmeraldError {
	return &LoadErrorInstance{
		Instance: rt.LoadError.New(),
		message:  msg,
	}
}
