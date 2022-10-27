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

	DefineSingletonMethod(ArgumentError, "new", exceptionNew(newArgumentError))
}

func newArgumentError(msg string) object.EmeraldError {
	return &ArgumentErrorInstance{
		Instance: ArgumentError.New(),
		message:  msg,
	}
}

func NewArgumentError(given int, expected int) object.EmeraldError {
	return newArgumentError(fmt.Sprintf("wrong number of arguments (given %d, expected %d)", given, expected))
}

func NewKeywordMissingArgumentError(keyword string) object.EmeraldError {
	return newArgumentError(fmt.Sprintf("missing keyword: :%s", keyword))
}
