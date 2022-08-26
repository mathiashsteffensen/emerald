package core

import (
	"emerald/object"
	"fmt"
)

var ArgumentError *object.Class

type ArgumentErrorInstance struct {
	*object.Instance
	message string
}

func (err *ArgumentErrorInstance) Message() string   { return err.message }
func (err *ArgumentErrorInstance) ClassName() string { return ArgumentError.Name }

func (err *ArgumentErrorInstance) Inspect() string {
	return fmt.Sprintf("#<%s: %s>", ArgumentError.Name, err.message)
}

func InitArgumentError() {
	ArgumentError = DefineClass(Object, "ArgumentError", StandardError)
	DefineMethod(ArgumentError, "new", exceptionNew(func(msg string) object.EmeraldError {
		return &ArgumentErrorInstance{
			Instance: ArgumentError.New(),
			message:  msg,
		}
	}), true)
}

func NewArgumentError(given int, expected int) object.EmeraldError {
	return &ArgumentErrorInstance{
		Instance: ArgumentError.New(),
		message:  fmt.Sprintf("wrong number of arguments (given %d, expected %d)", given, expected),
	}
}
