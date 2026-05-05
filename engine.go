package emerald

import (
	"emerald/bytecode"
	"emerald/compiler"
	"emerald/core"
	"emerald/object"
	"emerald/vm"
	"fmt"
)

type Engine struct {
	Runtime *core.Runtime
}

func New() *Engine {
	rt := core.NewRuntime()
	rt.Init()

	rt.CompileBlock = func(fileName string, content string) *bytecode.Bytecode {
		return compiler.Compile(fileName, content, rt)
	}

	return &Engine{
		Runtime: rt,
	}
}

func (e *Engine) Eval(content string) (object.EmeraldValue, error) {
	return e.EvalFile("(irb)", content)
}

func (e *Engine) EvalFile(fileName string, content string) (object.EmeraldValue, error) {
	bc := compiler.Compile(fileName, content, e.Runtime)

	machine := vm.New(fileName, bc, e.Runtime)
	machine.Run()

	globalException := e.Runtime.Heap.GetGlobalVariableString("$!")
	if !globalException.IsNil() && globalException != e.Runtime.NULL {
		return object.EmeraldValue{}, fmt.Errorf("exception: %s", globalException.Inspect())
	}

	return machine.LastPoppedStackElem(), nil
}
