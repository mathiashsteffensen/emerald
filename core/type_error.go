package core

import (
	"emerald/object"
	"fmt"
)

var TypeError *object.Class

func InitTypeError() {
	TypeError = DefineClass(Object, "TypeError", StandardError)
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
	return TypeError.Name
}

func NewNoConversionTypeError(expected string, actual string) object.EmeraldError {
	return NewTypeError(fmt.Sprintf("no implicit conversion of %s into %s", actual, expected))
}

func NewTypeError(msg string) object.EmeraldError {
	return &TypeErrorInstance{
		Instance: TypeError.New(),
		message:  msg,
	}
}
