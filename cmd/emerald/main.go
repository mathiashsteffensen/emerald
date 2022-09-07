package main

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/heap"
	"emerald/log"
	"emerald/object"
	"emerald/parser"
	"emerald/parser/lexer"
	"emerald/types"
	"emerald/vm"
	"flag"
	"fmt"
	"github.com/pkg/profile"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var profilingEnabled = flag.Bool("profile", false, "EM_DEBUG=1 emerald --profile lib/main.rb")
var logHeapUsage = flag.Bool("logHeapUsage", false, "EM_DEBUG=1 emerald --logHeapUsage lib/main.rb")

func main() {
	log.ExperimentalWarning()

	flag.Parse()

	if log.IsLevel(log.InternalDebugLevel) && *profilingEnabled {
		log.InternalDebug("Running with profiling enabled")
		defer profile.Start().Stop()
	}

	osArgs := types.NewSlice[string](os.Args[1:]...)
	fileIndex := osArgs.FindIndex(func(arg string) bool {
		return !strings.HasPrefix(arg, "--")
	})

	if fileIndex == nil {
		log.Fatal("No file to run")
	}

	file := osArgs.Value[*fileIndex]

	absFile, err := filepath.Abs(file)
	checkError("Failed to make path absolute?", err)

	bytes, err := os.ReadFile(absFile)
	checkError("error reading file", err)

	l := lexer.New(lexer.NewInput(absFile, string(bytes)))
	p := parser.New(l)
	program := p.ParseAST()

	if len(p.Errors()) != 0 {
		log.FatalF("parser error: %s\n", p.Errors()[0])
	}

	c := compiler.New()

	err = c.Compile(program)
	checkError("Compilation failed", err)

	argv := []object.EmeraldValue{}

	types.NewSlice[string](osArgs.Value[*fileIndex+1:]...).Each(func(arg string) {
		argv = append(argv, core.NewString(arg))
	})

	core.MainObject.NamespaceDefinitionSet("ARGV", core.NewArray(argv))

	defer recoverWithStacktrace()

	if *logHeapUsage {
		go logHeapUsageRoutine()
	}

	machine := vm.New(absFile, c.Bytecode())
	machine.Run()

	if exception := heap.GetGlobalVariableString("$!"); exception != nil {
		exception := exception.(object.EmeraldError)
		log.FatalF("%s: %s", exception.ClassName(), exception.Message())
	}

	if machine.StackTop() != nil {
		log.InternalDebug("StackTop was not nil")
	}

	log.Shutdown()
}

func logHeapUsageRoutine() {
	m := runtime.MemStats{}

	for {
		time.Sleep(200 * time.Millisecond)

		runtime.ReadMemStats(&m)
		heapAlloc := float64(m.HeapAlloc) / 1024 / 1024 // In MB

		log.DebugF("Heap size: %fMB", heapAlloc)
	}
}

func recoverWithStacktrace() {
	if r := recover(); r != nil {
		log.StackTrace(r)
	}
}

func checkError(msg string, err error) {
	if err != nil {
		fmt.Printf(msg+": %s\n", err.Error())
		os.Exit(1)
	}
}
