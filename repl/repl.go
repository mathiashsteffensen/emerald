package repl

import (
	"emerald"
	"emerald/compiler"
	"emerald/debug"
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

	var line string
	astNodes := []*ast.AST{}
	var buffer string
	var machine *vm.VM
	
	engine := emerald.New()
	comp := compiler.New(nil, engine.Runtime)
	lineCount := 1

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
				lineReader.SetPrompt(fmt.Sprintf(PROMPT_FMT, lineCount))
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

		oldInstructionsCount := len(comp.Bytecode().Instructions)

		comp.SetLexer(l)
		comp.Compile(program)

		bc := comp.Bytecode()

		if config.OutputBytecode {
			debug.InternalDebugF("Emerald bytecode: \n%s", bc.Instructions[oldInstructionsCount:])
			time.Sleep(50 * time.Millisecond)
		}

		if machine == nil {
			currentWorkingDir, err := os.Getwd()
			if err != nil {
				debug.Fatal(err.Error())
			}
			machine = vm.New(currentWorkingDir, bc, engine.Runtime)
		}

		newInstructions := bc.Instructions[oldInstructionsCount:]
		newDebugTokens := make(map[int]lexer.Token)
		for pos, token := range bc.DebugTokens {
			if pos >= oldInstructionsCount {
				newDebugTokens[pos-oldInstructionsCount] = token
			}
		}

		machine.RunIncremental(newInstructions, newDebugTokens, 0)

		if exception := engine.Runtime.Heap.GetGlobalVariableString("$!"); !exception.IsNil() && exception != engine.Runtime.NULL {
			printf("%s: %s\n", exception.Heap.(object.EmeraldError).ClassName(), exception.Heap.(object.EmeraldError).Message())
			engine.Runtime.Heap.SetGlobalVariableString("$!", object.EmeraldValue{})
			continue
		}

		evaluated := machine.LastPoppedStackElem()

		if !evaluated.IsNil() {
			evaluated = machine.Send(evaluated, "inspect", engine.Runtime.NULL, map[string]object.EmeraldValue{})
			print("=> " + evaluated.Inspect())
		}

		lineCount++
		lineReader.SetPrompt(fmt.Sprintf(PROMPT_FMT, lineCount))
	}

Exit:
	printf("\nSee you next time!")
}
