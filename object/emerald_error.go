package object

type EmeraldError interface {
	HeapObject
	Message() string
	ClassName() string
	// Backtrace is the call stack captured when the error was first raised, outermost frame last. nil until raised.
	Backtrace() []string
	SetBacktrace(backtrace []string)
}
