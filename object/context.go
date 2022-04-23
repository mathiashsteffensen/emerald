package object

type Context struct {
	Outer            *Context
	ExecutionTarget  EmeraldValue
	DefinitionTarget EmeraldValue
}
