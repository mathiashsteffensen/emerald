package core

import "emerald/object"

type ExceptionInstance struct {
	*object.Instance
	message string
}

func (err *ExceptionInstance) Message() string   { return err.message }
func (err *ExceptionInstance) ClassName() string { return "Exception" }

func (rt *Runtime) InitException() {
	rt.Exception = rt.DefineClass("Exception", rt.Object)

	rt.DefineSingletonMethod(rt.Exception, "new", rt.exceptionNew(rt.NewException))

	rt.DefineMethod(rt.Exception, "to_s", rt.exceptionToS())
}

func (rt *Runtime) NewException(msg string) object.EmeraldError {
	return &ExceptionInstance{
		Instance: rt.Exception.New(),
		message:  msg,
	}
}

func (rt *Runtime) exceptionToS() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString(ctx.Self.(object.EmeraldError).Message())
	}
}

func (rt *Runtime) exceptionNew(initializer func(msg string) object.EmeraldError) object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		var msg string

		if len(args) != 0 {
			msg = args[0].Inspect()
		}

		return initializer(msg)
	}
}
