package main

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/lexer"
	"emerald/parser"
	"emerald/vm"
	"fmt"
	"os"
)

func main() {
	file := os.Args[1]

	bytes, err := os.ReadFile(file)
	checkError("error reading file", err)

	l := lexer.New(lexer.NewInput(file, string(bytes)))
	p := parser.New(l)
	program := p.ParseAST()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Printf("parser error: %s\n", err)
		}
		os.Exit(1)
	}

	c := compiler.New(compiler.WithBuiltIns())

	err = c.Compile(program)
	checkError("Compilation failed", err)

	machine := vm.New(c.Bytecode())

	err = machine.Run()
	checkError("VM failed to execute bytecode", err)

	evaluated := machine.LastPoppedStackElem()

	if core.IsError(evaluated) {
		fmt.Println(evaluated.Inspect())
	}
}

func checkError(msg string, err error) {
	if err != nil {
		fmt.Printf(msg+": %s\n", err.Error())
		os.Exit(1)
	}
}
