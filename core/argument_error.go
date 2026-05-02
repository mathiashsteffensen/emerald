package core

import (
	"emerald/object"
	"fmt"
)

type ArgumentErrorInstance struct {
	*object.Instance
	message string
}

func (err *ArgumentErrorInstance) Message() string   { return err.message }
func (err *ArgumentErrorInstance) ClassName() string { return "ArgumentError" }

func (err *ArgumentErrorInstance) Inspect() string {
	return fmt.Sprintf("#<%s: %s>", "ArgumentError", err.message)
}

func (rt *Runtime) InitArgumentError() {
	rt.ArgumentError = rt.DefineClass("ArgumentError", rt.StandardError)

	rt.DefineSingletonMethod(rt.ArgumentError, "new", rt.exceptionNew(rt.newArgumentError))
}

func (rt *Runtime) newArgumentError(msg string) object.EmeraldError {
	return &ArgumentErrorInstance{
		Instance: rt.ArgumentError.New(),
		message:  msg,
	}
}

func (rt *Runtime) NewArgumentError(given int, expected int) object.EmeraldError {
	return rt.newArgumentError(fmt.Sprintf("wrong number of arguments (given %d, expected %d)", given, expected))
}

func (rt *Runtime) NewKeywordMissingArgumentError(keyword string) object.EmeraldError {
	return rt.newArgumentError(fmt.Sprintf("missing keyword: :%s", keyword))
}
