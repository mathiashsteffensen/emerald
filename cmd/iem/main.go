package main

import (
	"emerald/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
