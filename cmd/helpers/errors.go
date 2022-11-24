package helpers

import (
	"emerald/debug"
	"fmt"
	"os"
)

func RecoverWithStacktrace() {
	if r := recover(); r != nil {
		debug.StackTrace(r)
	}
}

func CheckError(msg string, err error) {
	if err != nil {
		fmt.Printf(msg+": %s\n", err.Error())
		os.Exit(1)
	}
}
