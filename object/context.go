package object

type Context struct {
	Outer *Context
	Self  EmeraldValue
}
