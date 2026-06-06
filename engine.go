package emerald

import (
	"emerald/bytecode"
	"emerald/compiler"
	"emerald/core"
	"emerald/object"
	"emerald/vm"
	"errors"
	"fmt"
)

type Engine struct {
	Runtime *core.Runtime
}

func New() *Engine {
	rt := core.NewRuntime()
	rt.Init()

	rt.CompileBlock = func(fileName string, content string) *bytecode.Bytecode {
		return compiler.CompileBlock(fileName, content, rt)
	}

	return &Engine{
		Runtime: rt,
	}
}

func (e *Engine) Eval(content string) (object.EmeraldValue, error) {
	return e.EvalFile("(irb)", content)
}

type EvalError struct {
	error
	Stacktrace []byte
}

func (e *Engine) EvalFile(fileName string, content string) (object.EmeraldValue, error) {
	bc := compiler.Compile(fileName, content, e.Runtime)

	machine := vm.New(fileName, bc, e.Runtime)
	machine.Run()

	globalException := e.Runtime.Heap.GetGlobalVariableString("$!")
	if !globalException.IsNil() && globalException != e.Runtime.NULL {
		if emErr, ok := globalException.Heap.(object.EmeraldError); ok {
			message := fmt.Sprintf("%s: %s", emErr.ClassName(), emErr.Message())
			return object.EmeraldValue{}, EvalError{error: errors.New(message), Stacktrace: []byte("")}
		}

		return object.EmeraldValue{}, fmt.Errorf("non-emerald-error exception: %s", globalException.Inspect())
	}

	return machine.LastPoppedStackElem(), nil
}
