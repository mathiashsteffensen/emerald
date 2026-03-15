package core

import (
	"emerald/object"
	"fmt"
)

var NameError *object.Class

func InitNameError() {
	NameError = DefineClass("NameError", StandardError)

	DefineSingletonMethod(NameError, "new", exceptionNew(NewNameError))
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
	return NameError.Name
}

func NewNameError(msg string) object.EmeraldError {
	return &NameErrorInstance{NameError.New(), msg}
}
