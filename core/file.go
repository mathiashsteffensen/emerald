package core

import (
	"emerald/object"
	"path/filepath"
)

func (rt *Runtime) InitFile() {
	rt.File = rt.DefineClass("File", rt.IO)

	rt.DefineSingletonMethod(rt.File, "absolute_path?", rt.fileIsAbsolutePath())
}

func (rt *Runtime) fileIsAbsolutePath() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		_, err := rt.EnforceArity(args, kwargs, 1, 1)
		if err != nil {
			return err
		}

		path, err := EnforceArgumentType[*StringInstance](rt, rt.String, args[0])
		if err != nil {
			return err
		}

		return rt.NativeBoolToBooleanObject(filepath.IsAbs(path.Value))
	}
}
