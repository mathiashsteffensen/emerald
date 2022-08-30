package repl

import (
	"emerald/compiler"
	"emerald/log"
	"emerald/parser"
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"emerald/vm"
	"fmt"
	"github.com/chzyer/readline"
	"io"
)

const PROMPT_FMT = "iem(main):%03d:0> "

type Config struct {
	OutputBytecode bool
	AstMode        bool
}

func Start(in io.ReadCloser, out io.Writer, config Config) {
	readline.SetHistoryPath("/tmp/iem.hst")

	buffer, err := readline.New(fmt.Sprintf(PROMPT_FMT, 1))
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize REPL buffer %s", err))
	}

	defer buffer.Close()

	buffer.Config.Stdin = in
	buffer.Config.Stdout = out
	buffer.Config.Stderr = out

	lineCount := 1

	var line string

	astNodes := []*ast.AST{}

	for {
		fmt.Fprintf(out, PROMPT_FMT, lineCount)

		line, err = buffer.Readline()
		if err != nil {
			if err.Error() == "Interrupt" {
				continue
			}

			switch err.Error() {
			case "Interrupt":
				continue
			default:
				fmt.Fprintf(out, "Error reading line %s\n", err)
				continue
			}
		}

		buffer.SaveHistory(line)

		if line == "quit" || line == "exit" {
			fmt.Fprintf(out, "See you next time!\n")
			break
		}

		l := lexer.New(lexer.NewInput("repl.rb", line))
		p := parser.New(l)
		program := p.ParseAST()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		if config.AstMode {
			astNodes = append(astNodes, program)
			for _, node := range astNodes {
				fmt.Fprintf(out, "%s\n", node.String())
			}
			continue
		}

		comp := compiler.New()

		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		code := comp.Bytecode()

		if config.OutputBytecode {
			log.InternalDebugF("Emerald bytecode: \n%s", code.Instructions[0:])
		}

		machine := vm.New(code)
		machine.Run()

		evaluated := machine.LastPoppedStackElem()

		if evaluated != nil {
			evaluated, _ = evaluated.SEND(machine.Context(), machine.Yield, "inspect", evaluated, nil)
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

		lineCount++
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
