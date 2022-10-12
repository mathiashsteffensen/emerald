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
		fmt.Fprintf(out, PROMPT_FMT, lineCount)

		line, err = lineReader.Readline()
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

		lineReader.SaveHistory(buffer + line)

		if line == "quit" || line == "exit" {
			fmt.Fprintf(out, "See you next time!\n")
			break
		}

		l := lexer.New(lexer.NewInput("repl.rb", buffer+line))
		p := parser.New(l)
		program := p.ParseAST()

		if len(p.Errors()) != 0 {
			errors := types.NewSlice(p.Errors()...)

			if errors.Includes("expected next token to be END, got EOF instead") ||
				errors.Includes("expected next token to be one of [RESCUE, ENSURE, END], got EOF instead") {
				buffer += line + "\n"
				lineReader.SetPrompt(lineReader.Config.Prompt + "	")
			} else {
				printParserErrors(out, p.Errors())
				lineReader.Config.Prompt = fmt.Sprintf(PROMPT_FMT, 1)
			}

			continue
		} else {
			buffer = ""
		}

		if config.AstMode {
			astNodes = append(astNodes, program)
			for _, node := range astNodes {
				fmt.Fprintf(out, "%s\n", node.String(0))
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
			log.FatalF("%s: %s", exception.ClassName(), exception.Message())
		}

		evaluated := machine.LastPoppedStackElem()

		if evaluated != nil {
			ctx := machine.Context()
			ctx.Self = evaluated
			evaluated = machine.Send(evaluated, "inspect", core.NULL)
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
