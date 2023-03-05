package core

import (
	"emerald/object"
	"fmt"
	"os"
	"path/filepath"
)

var Dir *object.Class

func InitDir() {
	Dir = DefineClass("Dir", Object)

	DefineSingletonMethod(Dir, "pwd", dirPwd())
	DefineSingletonMethod(Dir, "glob", dirGlob())
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

func dirGlob() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, emErr := EnforceArity(args, kwargs, 1, 2); emErr != nil {
			return emErr
		}

		globPath, emErr := EnforceArgumentType[*StringInstance](String, args[0])
		if emErr != nil {
			return emErr
		}

		paths, err := filepath.Glob(globPath.Value)
		if err != nil {
			Raise(NewRuntimeError(fmt.Sprintf("Failed to read glob path %s", err)))
			return NULL
		}

		result := NewArray([]object.EmeraldValue{})

		for _, path := range paths {
			result.Value = append(result.Value, NewString(path))
		}

		return result
	}
}
