package core

import "emerald/object"

var Exception *object.Class

type ExceptionInstance struct {
	*object.Instance
	message string
}

func (err *ExceptionInstance) Message() string   { return err.message }
func (err *ExceptionInstance) ClassName() string { return Exception.Name }

func InitException() {
	Exception = DefineClass(Object, "Exception", Object)
	
	DefineSingletonMethod(Exception, "new", exceptionNew(NewException))

	DefineMethod(Exception, "to_s", exceptionToS())
}

func NewException(msg string) object.EmeraldError {
	return &ExceptionInstance{
		Instance: Exception.New(),
		message:  msg,
	}
}

func exceptionToS() object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(self.(object.EmeraldError).Inspect())
	}
}

func exceptionNew(initializer func(msg string) object.EmeraldError) object.BuiltInMethod {
	return func(ctx *object.Context, self object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		var msg string

		if len(args) != 0 {
			msg = args[0].Inspect()
		}

		return initializer(msg)
	}
}
