package emerald

import (
	"emerald/bytecode"
	"emerald/compiler"
	"emerald/core"
	"emerald/object"
	"emerald/parser"
	"emerald/parser/lexer"
	"emerald/vm"
	"fmt"
)

type Engine struct {
	Runtime *core.Runtime
}

func New() *Engine {
	rt := core.NewRuntime()
	rt.Init()

	rt.Compile = func(fileName string, content string) *bytecode.Bytecode {
		return compiler.CompileToBytecode(fileName, content, rt)
	}

	return &Engine{
		Runtime: rt,
	}
}

func (e *Engine) Eval(content string) (object.EmeraldValue, error) {
	return e.EvalFile("(irb)", content)
}

func (e *Engine) EvalFile(fileName string, content string) (object.EmeraldValue, error) {
	l := lexer.New(lexer.NewInput(fileName, content))
	p := parser.New(l)
	ast := p.ParseAST()

	if len(p.Errors()) != 0 {
		return nil, fmt.Errorf("failed to parse source file %s\n\n%s", fileName, p.Errors()[0])
	}

	c := compiler.New(l, e.Runtime)
	c.Compile(ast)

	bc := c.Bytecode()
	bc.Instructions = append(bc.Instructions, byte(bytecode.OpReturn))

	machine := vm.New(fileName, bc, e.Runtime)
	machine.Run()

	globalException := e.Runtime.Heap.GetGlobalVariableString("$!")
	if globalException != nil && globalException != e.Runtime.NULL {
		return nil, fmt.Errorf("exception: %s", globalException.Inspect())
	}

	return machine.LastPoppedStackElem(), nil
}
