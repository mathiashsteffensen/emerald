package core

import "emerald/object"

var NoMethodError *object.Class

type NoMethodErrorInstance struct {
	*object.Instance
	message string
}

func (err *NoMethodErrorInstance) Message() string   { return err.message }
func (err *NoMethodErrorInstance) ClassName() string { return NoMethodError.Name }

func InitNoMethodError() {
	NoMethodError = DefineClass("NoMethodError", StandardError)

	DefineSingletonMethod(NoMethodError, "new", exceptionNew(NewNoMethodError))
}

func NewNoMethodError(msg string) object.EmeraldError {
	return &NoMethodErrorInstance{
		Instance: NoMethodError.New(),
		message:  msg,
	}
}
