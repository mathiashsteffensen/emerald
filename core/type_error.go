package core

import (
	"emerald/object"
	"fmt"
)

func (rt *Runtime) InitTypeError() {
	rt.TypeError = rt.DefineClass("TypeError", rt.StandardError)

	rt.DefineSingletonMethod(rt.TypeError, "new", rt.exceptionNew(rt.NewTypeError))
}

type TypeErrorInstance struct {
	*object.Instance
	message string
}

func (err *TypeErrorInstance) Inspect() string {
	return fmt.Sprintf("#<TypeError: %s>", err.message)
}

func (err *TypeErrorInstance) Message() string {
	return err.message
}

func (err *TypeErrorInstance) ClassName() string {
	return "TypeError"
}

func (rt *Runtime) NewNoConversionTypeError(expected string, actual string) object.EmeraldError {
	return rt.NewTypeError(fmt.Sprintf("no implicit conversion of %s into %s", actual, expected))
}

func (rt *Runtime) NewTypeError(msg string) object.EmeraldError {
	return &TypeErrorInstance{
		Instance: rt.TypeError.New(),
		message:  msg,
	}
}
