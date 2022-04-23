package core

import "emerald/object"

var Exception *object.Class

type ExceptionInstance struct {
	*object.Instance
	message string
}

func (err *ExceptionInstance) Message() string { return err.message }

func InitException() {
	Exception = object.NewClass(
		"Exception",
		Object,
		object.BuiltInMethodSet{
			"to_s": func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
				return NewString(target.(object.EmeraldError).Inspect())
			},
		},
		object.BuiltInMethodSet{
			"new": exceptionNew(NewException),
		},
	)
}

func NewException(msg string) object.EmeraldError {
	return &ExceptionInstance{
		Instance: Exception.New(),
		message:  msg,
	}
}

func exceptionNew(initializer func(msg string) object.EmeraldError) object.BuiltInMethod {
	return func(ctx *object.Context, target object.EmeraldValue, block object.EmeraldValue, yield object.YieldFunc, args ...object.EmeraldValue) object.EmeraldValue {
		var msg string

		if len(args) != 0 {
			msg = args[0].Inspect()
		}

		return initializer(msg)
	}
}
