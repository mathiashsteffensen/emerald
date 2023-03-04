package core

import (
	"emerald/object"
	"os"
)

var Dir *object.Class

func InitDir() {
	Dir = DefineClass("Dir", Object)

	DefineSingletonMethod(Dir, "pwd", dirPwd())
}

func dirPwd() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 0, 0); err != nil {
			return err
		}

		wd, err := os.Getwd()
		if err != nil {
			e := NewException(err.Error())
			Raise(e)
			return e
		}

		return NewString(wd)
	}
}
