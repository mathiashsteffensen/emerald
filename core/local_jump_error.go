package core

import (
	"emerald/object"
	"fmt"
)

func (rt *Runtime) InitLocalJumpError() {
	rt.LocalJumpError = rt.DefineClass("LocalJumpError", rt.StandardError)

	rt.defineNativeConstructor(rt.LocalJumpError, rt.exceptionNew(rt.NewLocalJumpError))
}

type LocalJumpErrorInstance struct {
	*object.Instance
	message string
}

func (err *LocalJumpErrorInstance) Inspect() string {
	return fmt.Sprintf("#<LocalJumpError: %s>", err.message)
}

func (err *LocalJumpErrorInstance) Message() string {
	return err.message
}

func (err *LocalJumpErrorInstance) ClassName() string {
	return "LocalJumpError"
}

func (rt *Runtime) NewLocalJumpError(msg string) object.EmeraldError {
	return &LocalJumpErrorInstance{rt.LocalJumpError.Heap.(*object.Class).New(), msg}
}
