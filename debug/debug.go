package debug

import (
	"emerald/types"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"strings"
	"time"
)

const EMERALD_VERSION = "0.0.1"

type LogLevel int

const (
	LogInternalDebugLevel LogLevel = iota
	LogDebugLevel
	LogWarnLevel
	LogFatalLevel
)

var currentLevel = LogDebugLevel
var trueEnvValues = types.NewSlice("true", "on", "1")
var IsTest bool
var BinaryDir string

func (l LogLevel) String() string {
	switch l {
	case LogInternalDebugLevel:
		return "INTERNAL"
	case LogDebugLevel:
		return "DEBUG"
	case LogWarnLevel:
		return "WARN"
	case LogFatalLevel:
		return "FATAL"
	}

	return ""
}

func IsLevel(level LogLevel) bool {
	return currentLevel == level
}

func init() {
	if setInternalDebug := os.Getenv("EM_DEBUG"); trueEnvValues.Includes(setInternalDebug) {
		currentLevel = LogInternalDebugLevel
	}

	isTestString := os.Getenv("EM_TEST")
	IsTest = trueEnvValues.Includes(isTestString)

	e, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		return
	}

	BinaryDir = path.Dir(e)

	go logRoutine()
}

type message struct {
	level  LogLevel
	format string
	args   []any
}

var msgChan = make(chan message, 50)
var doneChan = make(chan bool)

func Shutdown() {
	doneChan <- true
}

func writeToChan(level LogLevel, msg string) {
	write(level, msg)
}

func writeToChanF(level LogLevel, format string, args ...any) {
	writef(level, format, args...)
}

func InternalDebug(msg string) {
	writeToChan(LogInternalDebugLevel, msg)
}

func InternalDebugF(format string, args ...any) {
	writeToChanF(LogInternalDebugLevel, format, args...)
}

func Debug(msg string) {
	writeToChan(LogDebugLevel, msg)
}

func DebugF(format string, args ...any) {
	writeToChanF(LogDebugLevel, format, args...)
}

func Warn(msg string) {
	writeToChan(LogWarnLevel, msg)
}

func WarnF(format string, args ...any) {
	writeToChanF(LogWarnLevel, format, args...)
}

func Fatal(msg string) {
	writeToChan(LogFatalLevel, msg)
	os.Exit(1)
}

func FatalF(format string, args ...any) {
	writeToChanF(LogFatalLevel, format, args...)
	time.Sleep(200 * time.Millisecond)
	os.Exit(1)
}

func FatalBugF(format string, args ...any) {
	FatalF(
		"This is a bug in the Emerald toolchain, please report this issue at https://github.com/mathiashsteffensen/emerald/issues "+
			"with the error message and stack trace below.\n\n%s\n%s",
		fmt.Sprintf(format, args...),
		StackTrace(nil, false),
	)
}

func StackTrace(r any, log bool) string {
	goStack := string(debug.Stack())
	stackLines := strings.Join(strings.Split(goStack, "\n")[7:37], "\n")

	if log {
		FatalF("Emerald VM panicked %s:\n%s", r, stackLines)
	}

	return stackLines
}

func ExperimentalWarning() {
	WarnF(
		"Emerald VM version %s is experimental and not near complete. "+
			"You are guaranteed to encounter bugs. "+
			"When you do, feel free to report them at https://github.com/mathiashsteffensen/emerald/issues\n",
		EMERALD_VERSION,
	)
}

func write(level LogLevel, msg string) {
	if level >= currentLevel {
		fmt.Printf("[%s] %s\n", level, msg)
	}
}

func writef(level LogLevel, format string, args ...any) {
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
