package main

import (
	"emerald/log"
	"emerald/repl"
	"flag"
	"os"
)

var outputMode = flag.String("outmode", "inspect", "EM_DEBUG=1 iem -outmode=ast")

func main() {
	log.ExperimentalWarning()

	flag.Parse()

	repl.Start(os.Stdin, os.Stdout, repl.Config{OutputBytecode: *outputMode == "bytecode", AstMode: *outputMode == "ast"})
}
