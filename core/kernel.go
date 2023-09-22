package core

import (
	"emerald/debug"
	"emerald/object"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var Kernel *object.Module

var Compile func(fileName string, content string) []byte

func InitKernel() {
	Kernel = object.NewModule("Kernel", object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	DefineMethod(Kernel, "inspect", kernelInspect())
	DefineMethod(Kernel, "class", kernelClass())
	DefineMethod(Kernel, "kind_of?", kernelKindOf())
	DefineMethod(Kernel, "is_a?", kernelKindOf())
	DefineMethod(Kernel, "include", kernelInclude())

	definePrivateKernelMethod("raise", kernelRaise())
	definePrivateKernelMethod("require_relative", kernelRequireRelative())
	definePrivateKernelMethod("sleep", kernelSleep())
	definePrivateKernelMethod("puts", kernelPuts())
	definePrivateKernelMethod("print", kernelPrint())
}

func definePrivateKernelMethod(name string, method object.BuiltInMethod) {
	DefineMethod(Kernel, name, method, object.PRIVATE)
	DefineSingletonMethod(Kernel, name, method)
}

func kernelClass() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		class := ctx.Self.Class()

		for class.Type() != object.CLASS_VALUE {
			class = class.Super()
		}

		return class
	}
}

func kernelKindOf() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 1, 1); err != nil {
			return err
		}

		class := args[0]

		if class.Type() != object.CLASS_VALUE && class.Type() != object.MODULE_VALUE {
			Raise(NewTypeError("class or module required"))
			return NULL
		}

		for _, ancestor := range ctx.Self.Class().Ancestors() {
			if ancestor == args[0] {
				return TRUE
			}
		}

		return FALSE
	}
}

func kernelSleep() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		var sleepArg time.Duration

		switch arg := args[0].(type) {
		case *IntegerInstance:
			sleepArg = time.Duration(arg.Value) * time.Second
		case *FloatInstance:
			sleepArg = time.Duration(arg.Value) * time.Second
		}

		start := time.Now()
		time.Sleep(sleepArg)
		end := time.Now()
		slept := math.Round(end.Sub(start).Seconds())

		return NewInteger(int64(slept))
	}
}

func kernelPuts() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			val := Send(arg, "to_s", NULL, map[string]object.EmeraldValue{})

			if err := writeToStdout(fmt.Sprintf("%s\n", val.Inspect())); err != nil {
				return err
			}
		}

		return NULL
	}
}

func kernelPrint() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			val := Send(arg, "to_s", NULL, map[string]object.EmeraldValue{})

			if err := writeToStdout(val.Inspect()); err != nil {
				return err
			}
		}

		return NULL
	}
}

func kernelInclude() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 1, 255); err != nil {
			return err
		}

		for _, arg := range args {
			if arg == nil {
				continue
			}

			mod, ok := arg.(*object.Module)
			if !ok {
				Raise(NewTypeError(fmt.Sprintf("wrong argument type %s (expected Module)", arg.Class().Super().(*object.Class).Name)))
			}

			ctx.Self.Include(mod)
		}

		return ctx.Self
	}
}

func kernelRaise() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		args, err := EnforceArity(args, kwargs, 1, 2)
		if err != nil {
			return err
		}

		switch len(args) {
		case 1:
			Raise(NewRuntimeError(args[0].(*StringInstance).Value))
		case 2:
			exception := Send(args[0], "new", NULL, map[string]object.EmeraldValue{}, args[1])
			Raise(exception.(object.EmeraldError))
		}

		return NULL
	}
}

func kernelRequireRelative() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 1, 1); err != nil {
			return err
		}

		arg, emErr := EnforceArgumentType[*StringInstance](String, args[0])
		if emErr != nil {
			return emErr
		}

		path := arg.Value
		fileParts := strings.Split(ctx.File, "/")
		fileParts[len(fileParts)-1] = ""
		dir := filepath.Join(fileParts...)

		debug.InternalDebugF("Attempting to require %s from dir %s", path, dir)

		absoluteFilePath, err := filepath.Abs("/" + filepath.Join(dir, path))
		if err != nil {
			panic(err)
		}

		absolutePathStr := NewString(absoluteFilePath)

		// File has already been loaded
		if requiredFilesHash.Get(absolutePathStr) != NULL {
			debug.InternalDebugF("Kernel#require_relative - File %s is already loaded, skipping", absoluteFilePath)
			return FALSE
		}

		sourceContent, err := os.ReadFile(absoluteFilePath + ".rb")
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				Raise(NewLoadError(fmt.Sprintf("cannot load such file -- %s", absoluteFilePath)))
				return NULL
			}

			panic(err)
		}

		instructions := Compile(absoluteFilePath, string(sourceContent))

		debug.InternalDebugF("Kernel#require_relative - Successfully compiled file %s", absoluteFilePath)

		requiredBlock := object.NewClosedBlock(&object.Context{
			Outer: nil,
			File:  absoluteFilePath,
			Self:  MainObject,
			Block: NULL,
			Yield: ctx.Yield,
			BlockGiven: func() bool {
				return false
			},
		}, &object.Block{Instructions: instructions}, []object.EmeraldValue{}, "", object.PUBLIC)

		object.EvalBlock(requiredBlock, map[string]object.EmeraldValue{})

		requiredFilesHash.Set(absolutePathStr, TRUE)

		return TRUE
	}
}

func kernelInspect() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(ctx.Self.Inspect())
	}
}

var requiredFilesHash = NewHash()

func writeToStdout(str string) object.EmeraldError {
	_, err := os.Stdout.WriteString(str)
	if err != nil {
		return raiseStdoutWriteFailed()
	}

	return nil
}

func raiseStdoutWriteFailed() object.EmeraldError {
	err := NewException("Failed to write to STDOUT")
	Raise(err)
	return err
}
