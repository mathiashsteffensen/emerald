package repl

import (
	"bufio"
	"emerald/compiler"
	"emerald/lexer"
	"emerald/object"
	"emerald/parser"
	"emerald/vm"
	"fmt"
	"io"
)

const PROMPT_FMT = "iem(main):%03d:0> "

func Start(in io.Reader, out io.Writer) {
	lineCount := 0
	scanner := bufio.NewScanner(in)

	constants := []object.EmeraldValue{}
	globals := make([]object.EmeraldValue, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		lineCount++

		fmt.Fprintf(out, PROMPT_FMT, lineCount)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "quit" {
			fmt.Fprintf(out, "See you next time!\n")
			return
		}

		l := lexer.New(lexer.NewInput("repl.rb", line))
		p := parser.New(l)
		program := p.ParseAST()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.New(compiler.WithState(symbolTable, constants), compiler.WithBuiltIns())

		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.New(code, vm.WithGlobalsStore(globals))

		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		evaluated := machine.LastPoppedStackElem()

		if evaluated != nil {
			evaluated, _ = evaluated.SEND(machine.Context(), machine.EvalBlock, "inspect", evaluated, nil)
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
