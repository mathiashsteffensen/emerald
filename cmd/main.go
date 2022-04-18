package main

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/lexer"
	"emerald/parser"
	"emerald/vm"
	"fmt"
	"io"
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
		return
	}

	c := compiler.New()

	err = c.Compile(program)
	checkError("Compilation failed", err)

	machine := vm.New(c.Bytecode())

	err = machine.Run()
	checkError("VM failed to execute bytecode", err)

	evaluated := machine.LastPoppedStackElem()

	if evaluated != nil {
		if evaluated.RespondsTo("to_s", evaluated) {
			evaluated, err = evaluated.SEND(nil, "to_s", evaluated, nil)
			if err != nil {
				evaluated = core.NewStandardError(err.Error())
			}
		}

		io.WriteString(os.Stdout, evaluated.Inspect())
		io.WriteString(os.Stdout, "\n")
	}
}

func checkError(msg string, err error) {
	if err != nil {
		fmt.Printf(msg+": %s\n", err.Error())
		os.Exit(1)
	}
}
