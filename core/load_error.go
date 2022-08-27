package core

import (
	"emerald/object"
	"fmt"
)

var LoadError *object.Class

func InitLoadError() {
	LoadError = DefineClass(Object, "LoadError", StandardError)

	DefineSingletonMethod(LoadError, "new", exceptionNew(NewLoadError))
}

type LoadErrorInstance struct {
	*object.Instance
	message string
}

func (err *LoadErrorInstance) Inspect() string {
	return fmt.Sprintf("#<LoadError: %s>", err.message)
}

func (err *LoadErrorInstance) Message() string {
	return err.message
}

func (err *LoadErrorInstance) ClassName() string {
	return LoadError.Name
}

func NewLoadError(msg string) object.EmeraldError {
	return &LoadErrorInstance{
		Instance: LoadError.New(),
		message:  msg,
	}
}
