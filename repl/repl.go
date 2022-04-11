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

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	constants := []object.EmeraldValue{}
	globals := make([]object.EmeraldValue, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		fmt.Fprint(out, PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(lexer.NewInput("repl.rb", line))
		p := parser.New(l)
		program := p.ParseAST()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.NewWithState(symbolTable, constants)

		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := vm.NewWithGlobalsStore(code, globals)

		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}

		evaluated := machine.LastPoppedStackElem()

		if evaluated != nil {
			if evaluated.RespondsTo("to_s", evaluated) {
				io.WriteString(
					out,
					evaluated.
						SEND(
							func(block *object.Block, args ...object.EmeraldValue) object.EmeraldValue {
								return object.NULL
							},
							"to_s",
							evaluated,
							nil,
						).
						Inspect(),
				)
			} else {
				io.WriteString(out, evaluated.Inspect())
			}

			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
