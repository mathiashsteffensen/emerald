package object

type Context struct {
	Outer      *Context
	File       string
	Self       EmeraldValue
	BreakValue EmeraldValue
	Block      EmeraldValue
	Yield      func(args ...EmeraldValue) EmeraldValue
	BlockGiven func() bool
}
