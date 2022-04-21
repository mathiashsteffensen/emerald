package core

import (
	"emerald/object"
	"fmt"
)

var ArgumentError *object.Class

type ArgumentErrorInstance struct {
	*StandardErrorInstance
}

func init() {
	ArgumentError = object.NewClass("ArgumentError", StandardError, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
}

func NewArgumentError(given string, expected string) *ArgumentErrorInstance {
	return &ArgumentErrorInstance{
		StandardErrorInstance: &StandardErrorInstance{
			Instance: ArgumentError.New(),
			message:  fmt.Sprintf("wrong number of arguments (given %s, expected %s)", given, expected),
		},
	}
}
