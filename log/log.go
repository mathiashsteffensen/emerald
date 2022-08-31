package log

import (
	"emerald/types"
	"fmt"
	"os"
)

type Level int

const (
	InternalDebugLevel Level = iota
	// DebugLevel
	WarnLevel
	FatalLevel
)

var currentLevel = WarnLevel
var trueEnvValues = types.NewSlice("true", "on", "1")

func (l Level) String() string {
	switch l {
	case InternalDebugLevel:
		return "INTERNAL"
	case WarnLevel:
		return "WARN"
	case FatalLevel:
		return "FATAL"
	}

	return ""
}

func IsLevel(level Level) bool {
	return currentLevel == level
}

func init() {
	if setInternalDebug := os.Getenv("EM_DEBUG"); trueEnvValues.Includes(setInternalDebug) {
		currentLevel = InternalDebugLevel
	}

	go logRoutine()
}

type message struct {
	level  Level
	format string
	args   []any
}

var msgChan = make(chan message, 50)
var doneChan = make(chan bool)

func Shutdown() {
	doneChan <- true
}

func writeToChan(level Level, msg string) {
	msgChan <- message{
		level:  level,
		format: msg,
	}
}

func writeToChanF(level Level, format string, args ...any) {
	msgChan <- message{
		level:  level,
		format: format,
		args:   args,
	}
}

func InternalDebug(msg string) {
	writeToChan(InternalDebugLevel, msg)
}

func InternalDebugF(format string, args ...any) {
	writeToChanF(InternalDebugLevel, format, args...)
}

func Warn(msg string) {
	writeToChan(WarnLevel, msg)
}

func WarnF(format string, args ...any) {
	writeToChanF(WarnLevel, format, args...)
}

func Fatal(msg string) {
	writeToChan(FatalLevel, msg)
	os.Exit(1)
}

func FatalF(format string, args ...any) {
	writeToChanF(WarnLevel, format, args...)
	os.Exit(1)
}

func ExperimentalWarning() {
	Warn(
		"The Emerald VM is experimental and not near complete. " +
			"You are guaranteed to encounter bugs. " +
			"When you do, feel free to report them at https://github.com/mathiashsteffensen/emerald/issues\n",
	)
}

func write(level Level, msg string) {
	if level >= currentLevel {
		fmt.Printf("[%s] %s\n", level.String(), msg)
	}
}

func writef(level Level, format string, args ...any) {
	if level >= currentLevel {
		write(level, fmt.Sprintf(format, args...))
	}
}

func logRoutine() {
	for {
		select {
		case msg := <-msgChan:
			if msg.args == nil {
				write(msg.level, msg.format)
			} else {
				writef(msg.level, msg.format, msg.args...)
			}
		case <-doneChan:
			return
		}
	}
}
