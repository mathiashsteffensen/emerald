package core

import (
	"emerald/object"
	"path/filepath"
)

var File *object.Class

func InitFile() {
	File = DefineClass("File", IO)

	DefineSingletonMethod(File, "absolute_path?", fileIsAbsolutePath())
}

func fileIsAbsolutePath() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		_, err := EnforceArity(args, kwargs, 1, 1)
		if err != nil {
			return err
		}

		path, err := EnforceArgumentType[*StringInstance](String, args[0])
		if err != nil {
			return err
		}

		return NativeBoolToBooleanObject(filepath.IsAbs(path.Value))
	}
}
