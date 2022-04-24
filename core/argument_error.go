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

func (err *ArgumentErrorInstance) Message() string { return err.message }

func (err *ArgumentErrorInstance) Inspect() string {
	return fmt.Sprintf("#<%s: %s>", ArgumentError.Name, err.message)
}

func InitArgumentError() {
	ArgumentError = object.NewClass(
		"ArgumentError",
		StandardError,
		StandardError.Class(),
		object.BuiltInMethodSet{},
		object.BuiltInMethodSet{
			"new": exceptionNew(func(msg string) object.EmeraldError {
				return &ArgumentErrorInstance{
					Instance: ArgumentError.New(),
					message:  msg,
				}
			}),
		},
	)
}

func NewArgumentError(given string, expected string) object.EmeraldError {
	return &ArgumentErrorInstance{
		Instance: ArgumentError.New(),
		message:  fmt.Sprintf("wrong number of arguments (given %s, expected %s)", given, expected),
	}
}
