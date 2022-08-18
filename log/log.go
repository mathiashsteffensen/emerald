package log

import (
	"fmt"
)

type Level int

const (
	InternalDebugLevel Level = iota
	// Debug
	// Warn
)

var currentLevel = InternalDebugLevel

func write(level Level, msg string) {
	if level >= currentLevel {
		fmt.Println(msg)
	}
}

func writef(level Level, format string, args ...any) {
	write(level, fmt.Sprintf(format, args...))
}

func InternalDebug(msg string) {
	write(InternalDebugLevel, msg)
}

func InternalDebugF(format string, args ...any) {
	writef(InternalDebugLevel, format, args...)
}
