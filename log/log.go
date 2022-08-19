package log

import (
	"fmt"
)

type Level int

const (
	InternalDebugLevel Level = iota
	// DebugLevel
	WarnLevel
)

var currentLevel = InternalDebugLevel

func (l Level) String() string {
	switch l {
	case InternalDebugLevel:
		return "INTERNAL"
	case WarnLevel:
		return "WARN"
	}

	return ""
}

func write(level Level, msg string) {
	if level >= currentLevel {
		fmt.Printf("[%s] %s\n", level.String(), msg)
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

func Warn(msg string) {
	write(WarnLevel, msg)
}

func WarnF(format string, args ...any) {
	writef(WarnLevel, format, args...)
}

func ExperimentalWarning() {
	Warn(
		"The Emerald VM is experimental and not near complete. " +
			"You are guaranteed to encounter bugs. " +
			"When you do, feel free to report them at https://github.com/mathiashsteffensen/emerald/issues\n",
	)
}
