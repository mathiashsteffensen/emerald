package core

import (
	"emerald/object"
	"fmt"
)

var StandardError *object.Class

func InitStandardError() {
	StandardError = DefineClass("StandardError", Exception)

	DefineSingletonMethod(StandardError, "new", exceptionNew(NewStandardError))
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
	return StandardError.Name
}

func NewStandardError(msg string) object.EmeraldError {
	return &StandardErrorInstance{StandardError.New(), msg}
}
