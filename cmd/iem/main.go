package main

import (
	"emerald/log"
	"emerald/repl"
	"os"
)

func main() {
	log.ExperimentalWarning()

	repl.Start(os.Stdin, os.Stdout)
}
