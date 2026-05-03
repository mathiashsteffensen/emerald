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

func (rt *Runtime) InitKernel() {
	rt.Kernel = object.NewModule("Kernel", object.BuiltInMethodSet{}, object.BuiltInMethodSet{})

	rt.DefineMethod(rt.Kernel, "inspect", rt.kernelInspect())
	rt.DefineMethod(rt.Kernel, "class", rt.kernelClass())
	rt.DefineMethod(rt.Kernel, "kind_of?", rt.kernelKindOf())
	rt.DefineMethod(rt.Kernel, "is_a?", rt.kernelKindOf())
	rt.DefineMethod(rt.Kernel, "include", rt.kernelInclude())

	rt.definePrivateKernelMethod("raise", rt.kernelRaise())
	rt.definePrivateKernelMethod("require_relative", rt.kernelRequireRelative())
	rt.definePrivateKernelMethod("sleep", rt.kernelSleep())
	rt.definePrivateKernelMethod("puts", rt.kernelPuts())
	rt.definePrivateKernelMethod("print", rt.kernelPrint())
}

func (rt *Runtime) definePrivateKernelMethod(name string, method object.BuiltInMethod) {
	rt.DefineMethod(rt.Kernel, name, method, object.PRIVATE)
	rt.DefineSingletonMethod(rt.Kernel, name, method)
}

func (rt *Runtime) kernelClass() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		class := ctx.Self.Class()

		for class.Type() != object.CLASS_VALUE {
			class = class.Super()
		}

		return class
	}
}

func (rt *Runtime) kernelKindOf() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return err
		}

		class := args[0]

		if class.Type() != object.CLASS_VALUE && class.Type() != object.MODULE_VALUE {
			rt.Raise(rt.NewTypeError("class or module required"))
			return rt.NULL
		}

		for _, ancestor := range ctx.Self.Class().Ancestors() {
			if ancestor == args[0] {
				return rt.TRUE
			}
		}

		return rt.FALSE
	}
}

func (rt *Runtime) kernelSleep() object.BuiltInMethod {
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

		return rt.NewInteger(int64(slept))
	}
}

func (rt *Runtime) kernelPuts() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			val := rt.Send(arg, "to_s", rt.NULL, map[string]object.EmeraldValue{})

			if err := rt.writeToStdout(fmt.Sprintf("%s\n", val.Inspect())); err != nil {
				return err
			}
		}

		return rt.NULL
	}
}

func (rt *Runtime) kernelPrint() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			val := rt.Send(arg, "to_s", rt.NULL, map[string]object.EmeraldValue{})

			if err := rt.writeToStdout(val.Inspect()); err != nil {
				return err
			}
		}

		return rt.NULL
	}
}

func (rt *Runtime) kernelInclude() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 255); err != nil {
			return err
		}

		for _, arg := range args {
			if arg == nil {
				continue
			}

			mod, ok := arg.(*object.Module)
			if !ok {
				rt.Raise(rt.NewTypeError(fmt.Sprintf("wrong argument type %s (expected Module)", arg.Class().Super().(*object.Class).Name)))
			}

			ctx.Self.Include(mod)
		}

		return ctx.Self
	}
}

func (rt *Runtime) kernelRaise() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		args, err := rt.EnforceArity(args, kwargs, 1, 2)
		if err != nil {
			return err
		}

		switch len(args) {
		case 1:
			rt.Raise(rt.NewRuntimeError(args[0].(*StringInstance).Value))
		case 2:
			exception := rt.Send(args[0], "new", rt.NULL, map[string]object.EmeraldValue{}, args[1])
			rt.Raise(exception.(object.EmeraldError))
		}

		return rt.NULL
	}
}

func (rt *Runtime) kernelRequireRelative() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 1, 1); err != nil {
			return err
		}

		arg, emErr := EnforceArgumentType[*StringInstance](rt, rt.String, args[0])
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

		if rt.RequiredFilesHash == nil {
			panic("RequiredFilesHash is nil!")
		}

		absolutePathStr := rt.NewString(absoluteFilePath)

		// rt.File has already been loaded
		if rt.RequiredFilesHash.(*HashInstance).Get(absolutePathStr) != nil {
			debug.InternalDebugF("rt.Kernel#require_relative - rt.File %s is already loaded, skipping", absoluteFilePath)
			return rt.FALSE
		}

		sourceContent, err := os.ReadFile(absoluteFilePath + ".rb")
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				rt.Raise(rt.NewLoadError(fmt.Sprintf("cannot load such file -- %s", absoluteFilePath)))
				return rt.NULL
			}

			panic(err)
		}

		bytecode := rt.Compile(absoluteFilePath, string(sourceContent))

		debug.InternalDebugF("rt.Kernel#require_relative - Successfully compiled file %s", absoluteFilePath)

		requiredBlock := object.NewClosedBlock(&object.Context{
			Outer: nil,
			File:  absoluteFilePath,
			Self:  rt.MainObject,
			Block: rt.NULL,
			Yield: ctx.Yield,
			BlockGiven: func() bool {
				return false
			},
		}, &object.Block{Bytecode: *bytecode}, []object.EmeraldValue{}, "", object.PUBLIC)

		rt.EvalBlock(requiredBlock, map[string]object.EmeraldValue{})

		rt.RequiredFilesHash.(*HashInstance).Set(absolutePathStr, rt.TRUE)

		return rt.TRUE
	}
}

func (rt *Runtime) kernelInspect() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		return rt.NewString(ctx.Self.Inspect())
	}
}

func (rt *Runtime) writeToStdout(str string) object.EmeraldError {
	_, err := os.Stdout.WriteString(str)
	if err != nil {
		return rt.raiseStdoutWriteFailed()
	}

	return nil
}

func (rt *Runtime) raiseStdoutWriteFailed() object.EmeraldError {
	err := rt.NewException("Failed to write to STDOUT")
	rt.Raise(err)
	return err
}
