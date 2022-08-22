package main

import (
	"emerald/log"
	"emerald/repl"
	"flag"
	"os"
)

var outputBytecode = flag.Bool("output-bytecode", false, "EM_DEBUG=true iem --output-bytecode")

func main() {
	log.ExperimentalWarning()

	flag.Parse()

	repl.Start(os.Stdin, os.Stdout, repl.Config{OutputBytecode: *outputBytecode})
}
