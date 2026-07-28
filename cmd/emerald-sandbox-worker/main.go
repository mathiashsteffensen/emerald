package main

import (
	"emerald"
	"fmt"
	"os"
)

func main() {
	if err := emerald.ServeSandboxWorker(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
