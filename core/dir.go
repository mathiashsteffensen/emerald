package core

import (
	"emerald/object"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
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

var globListRegexp = regexp.MustCompile(`{(\S*)}`)

func raiseFailedToReadGlobPathError(err error) {
	Raise(NewRuntimeError(fmt.Sprintf("Failed to read glob path %s", err)))
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

		var paths []string

		for _, match := range globListRegexp.FindAllString(globPath.Value, -1) {
			values := strings.Split(match[1:len(match)-1], ",")

			for _, value := range values {
				p, err := filepath.Glob(
					globListRegexp.ReplaceAllString(globPath.Value, value),
				)
				if err != nil {
					raiseFailedToReadGlobPathError(err)
					return NULL
				}
				paths = append(paths, p...)
			}
		}

		p, err := filepath.Glob(globPath.Value)
		if err != nil {
			raiseFailedToReadGlobPathError(err)
			return NULL
		}
		paths = append(paths, p...)
		sort.Strings(paths)

		result := NewArray([]object.EmeraldValue{})

		for _, path := range paths {
			result.Value = append(result.Value, NewString(path))
		}

		return result
	}
}
