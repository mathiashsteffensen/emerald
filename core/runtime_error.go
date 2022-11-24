package core

import (
	"emerald/object"
	"fmt"
)

var RuntimeError *object.Class

func InitRuntimeError() {
	RuntimeError = DefineClass("RuntimeError", StandardError)

	DefineSingletonMethod(RuntimeError, "new", exceptionNew(NewRuntimeError))
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
	return RuntimeError.Name
}

func NewRuntimeError(msg string) object.EmeraldError {
	return &RuntimeErrorInstance{RuntimeError.New(), msg}
}
