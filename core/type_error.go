package core

import (
	"emerald/object"
	"fmt"
)

var TypeError *object.Class

func InitTypeError() {
	TypeError = object.NewClass("TypeError", StandardError, StandardError.Class(), object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
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

func NewTypeError(expected string, actual string) object.EmeraldError {
	return &TypeErrorInstance{
		Instance: TypeError.New(),
		message:  fmt.Sprintf("no implicit conversion of %s into %s", actual, expected),
	}
}
