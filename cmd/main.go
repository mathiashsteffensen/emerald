package main

import (
	"emerald/evaluator"
	"emerald/lexer"
	"emerald/object"
	"emerald/parser"
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

	env := object.NewEnvironment()
	ctx := object.ExecutionContext{Target: object.Object, IsStatic: true}
	evaluated := evaluator.Eval(ctx, program, env)

	if evaluated != nil {
		if evaluated.RespondsTo("to_s", evaluated) {
			io.WriteString(
				os.Stdout,
				evaluated.
					SEND(
						evaluator.Eval,
						evaluator.Yield(ctx),
						"to_s",
						evaluated,
						nil,
					).
					Inspect(),
			)
		} else {
			io.WriteString(os.Stdout, evaluated.Inspect())
		}

		io.WriteString(os.Stdout, "\n")
	}
}

func checkError(msg string, err error) {
	if err != nil {
		fmt.Printf(msg+": %s\n", err.Error())
		os.Exit(1)
	}
}
