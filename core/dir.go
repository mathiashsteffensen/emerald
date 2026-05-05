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

func (rt *Runtime) InitDir() {
	rt.Dir = rt.DefineClass("Dir", rt.Object)

	rt.DefineSingletonMethod(rt.Dir, "pwd", rt.dirPwd())
	rt.DefineSingletonMethod(rt.Dir, "glob", rt.dirGlob())
}

func (rt *Runtime) dirPwd() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 0, 0); err != nil {
			return object.NewHeapObject(err)
		}

		wd, err := os.Getwd()
		if err != nil {
			e := rt.NewException(err.Error())
			rt.Raise(e)
			return object.NewHeapObject(e)
		}

		return rt.NewString(wd)
	}
}

var globListRegexp = regexp.MustCompile(`{(\S*)}`)

func (rt *Runtime) raiseFailedToReadGlobPathError(err error) {
	rt.Raise(rt.NewRuntimeError(fmt.Sprintf("Failed to read glob path %s", err)))
}

func (rt *Runtime) dirGlob() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, emErr := rt.EnforceArity(args, kwargs, 1, 2); emErr != nil {
			return object.NewHeapObject(emErr)
		}

		globPath, emErr := EnforceArgumentType[*StringInstance](rt, rt.String, args[0])
		if emErr != nil {
			return object.NewHeapObject(emErr)
		}

		var paths []string

		for _, match := range globListRegexp.FindAllString(globPath.Value, -1) {
			values := strings.Split(match[1:len(match)-1], ",")

			for _, value := range values {
				p, err := filepath.Glob(
					globListRegexp.ReplaceAllString(globPath.Value, value),
				)
				if err != nil {
					rt.raiseFailedToReadGlobPathError(err)
					return rt.NULL
				}
				paths = append(paths, p...)
			}
		}

		p, err := filepath.Glob(globPath.Value)
		if err != nil {
			rt.raiseFailedToReadGlobPathError(err)
			return rt.NULL
		}
		paths = append(paths, p...)
		sort.Strings(paths)

		result := rt.NewArray([]object.EmeraldValue{})

		for _, path := range paths {
			result.Heap.(*ArrayInstance).Value = append(result.Heap.(*ArrayInstance).Value, rt.NewString(path))
		}

		return result
	}
}
