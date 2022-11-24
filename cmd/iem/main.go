package main

import (
	"emerald/debug"
	"emerald/repl"
	"flag"
	"os"
)

var outputMode = flag.String("outmode", "inspect", "EM_DEBUG=1 iem -outmode=ast")

func main() {
	debug.ExperimentalWarning()

	flag.Parse()

	repl.Start(os.Stdin, os.Stdout, repl.Config{OutputBytecode: *outputMode == "bytecode", AstMode: *outputMode == "ast"})
}
