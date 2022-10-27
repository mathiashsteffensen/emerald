package repl

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/heap"
	"emerald/log"
	"emerald/object"
	"emerald/parser"
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"emerald/types"
	"emerald/vm"
	"fmt"
	"github.com/chzyer/readline"
	"io"
	"os"
	"time"
)

const PROMPT_FMT = "iem(main):%03d:0> "

type Config struct {
	OutputBytecode bool
	AstMode        bool
}

func Start(in io.ReadCloser, out io.Writer, config Config) {
	print := func(str string) {
		_, err := io.WriteString(out, str)
		if err != nil {
			panic(err)
		}
		_, err = io.WriteString(out, "\n")
		if err != nil {
			panic(err)
		}
	}

	printf := func(frmt string, args ...any) {
		print(fmt.Sprintf(frmt, args...))
	}

	readline.SetHistoryPath("/tmp/iem.hst")

	lineReader, err := readline.New(fmt.Sprintf(PROMPT_FMT, 1))
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize REPL lineReader %s", err))
	}

	defer lineReader.Close()

	lineReader.Config.Stdin = in
	lineReader.Config.Stdout = out
	lineReader.Config.Stderr = out

	lineCount := 1

	var line string

	astNodes := []*ast.AST{}

	var buffer string

	for {
		line, err = lineReader.Readline()
		if err != nil {
			if err.Error() == "Interrupt" {
				continue
			}

			switch err.Error() {
			case "Interrupt":
				continue
			case "EOF":
				goto Exit
			default:
				printf("Error reading line %s\n", err)
				continue
			}
		}

		lineReader.SaveHistory(buffer + line)

		if line == "quit" || line == "exit" {
			goto Exit
		}

		l := lexer.New(lexer.NewInput("repl.rb", buffer+line))
		p := parser.New(l)
		program := p.ParseAST()

		if len(p.Errors()) != 0 {
			errors := types.NewSlice(p.Errors()...)

			if errors.Includes("syntax error, unexpected end-of-input") {
				buffer += line + "\n"
				lineReader.SetPrompt(lineReader.Config.Prompt + "	")
			} else {
				for _, msg := range p.Errors() {
					print("\t" + msg)
				}
				lineReader.Config.Prompt = fmt.Sprintf(PROMPT_FMT, 1)
			}

			continue
		} else {
			buffer = ""
			lineReader.SetPrompt(fmt.Sprintf(PROMPT_FMT, 1))
		}

		if config.AstMode {
			astNodes = append(astNodes, program)
			for _, node := range astNodes {
				printf("%s\n", node.String(0))
			}
			continue
		}

		comp := compiler.New()
		comp.Compile(program)

		code := comp.Bytecode()

		if config.OutputBytecode {
			log.InternalDebugF("Emerald bytecode: \n%s", code.Instructions[0:])
			time.Sleep(50 * time.Millisecond)
		}

		currentWorkingDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err.Error())
		}

		machine := vm.New(currentWorkingDir, code)
		machine.Run()

		if exception := heap.GetGlobalVariableString("$!"); exception != nil {
			exception := exception.(object.EmeraldError)
			printf("%s: %s\n", exception.ClassName(), exception.Message())
			heap.SetGlobalVariableString("$!", nil)
			continue
		}

		evaluated := machine.LastPoppedStackElem()

		if evaluated != nil {
			ctx := machine.Context()
			ctx.Self = evaluated
			evaluated = machine.Send(evaluated, "inspect", core.NULL)
			print(evaluated.Inspect())
		}

		lineCount++
	}

Exit:
	printf("\nSee you next time!")
}
