package main

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/log"
	"emerald/parser"
	"emerald/parser/lexer"
	"emerald/types"
	"emerald/vm"
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"os"
	"strings"
)

var profilingEnabled = flag.Bool("profile", false, "EM_DEBUG=1 emerald --profile lib/main.rb")

func main() {
	log.ExperimentalWarning()

	flag.Parse()

	if log.IsLevel(log.InternalDebugLevel) && *profilingEnabled {
		log.InternalDebug("Running with profiling enabled")
		defer profile.Start().Stop()
	}

	args := types.NewSlice(os.Args[1:]...)

	file := args.Find(func(arg string) bool {
		return !strings.HasPrefix(arg, "--")
	})

	if file == nil {
		log.Fatal("No file to run")
	}

	bytes, err := os.ReadFile(*file)
	checkError("error reading file", err)

	l := lexer.New(lexer.NewInput(*file, string(bytes)))
	p := parser.New(l)
	program := p.ParseAST()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Printf("parser error: %s\n", err)
		}
		os.Exit(1)
	}

	c := compiler.New()

	err = c.Compile(program)
	checkError("Compilation failed", err)

	machine := vm.New(c.Bytecode())
	machine.Run()

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
