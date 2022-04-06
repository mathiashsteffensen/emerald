package repl

import (
	"bufio"
	"emerald/evaluator"
	"emerald/lexer"
	"emerald/object"
	"emerald/parser"
	"fmt"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "quit" {
			return
		}

		l := lexer.New(lexer.NewInput("interactive.rb", line))
		p := parser.New(l)
		program := p.ParseAST()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(object.Object, program, env)

		if evaluated != nil {
			if evaluated.RespondsTo("to_s", evaluated) {
				io.WriteString(
					out,
					evaluated.
						SEND(
							evaluator.Eval,
							env,
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
