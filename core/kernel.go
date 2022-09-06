package core

import (
	"emerald/log"
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
	Kernel = DefineModule(nil, "Kernel")

	DefineMethod(Kernel, "inspect", kernelInspect())
	DefineMethod(Kernel, "raise", kernelRaise())
	DefineMethod(Kernel, "require_relative", kernelRequireRelative())
	DefineMethod(Kernel, "class", kernelClass())
	DefineMethod(Kernel, "kind_of?", kernelKindOf())
	DefineMethod(Kernel, "is_a?", kernelKindOf())
	DefineMethod(Kernel, "include", kernelInclude())
	DefineMethod(Kernel, "sleep", kernelSleep())

	// Should be made private when that function has been implemented
	DefineMethod(Kernel, "puts", kernelPuts())
}

func kernelClass() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		class := ctx.Self.Class()

		for class.Type() != object.CLASS_VALUE {
			class = class.Super()
		}

		return class
	}
}

func kernelKindOf() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		if len(args) != 1 {
			return NewArgumentError(len(args), 1)
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
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
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
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		for _, arg := range args {
			val := Send(arg, "to_s", NULL)

			if err := writeToStdout(val.Inspect()); err != nil {
				return err
			}
		}

		return NULL
	}
}

func kernelInclude() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		if len(args) == 0 {
			panic(errors.New("wrong number of arguments (given 0, expected 1+)"))
		}

		for _, arg := range args {
			if arg == nil {
				continue
			}

			mod, ok := arg.(*object.Module)
			if !ok {
				panic(fmt.Errorf("wrong argument type %s (expected Module)", arg.Inspect()))
			}

			ctx.Self.Include(mod)
		}

		return ctx.Self
	}
}

func kernelRaise() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		args, err := EnforceArity(args, 1, 2)
		if err != nil {
			return err
		}

		switch len(args) {
		case 1:
			Raise(NewRuntimeError(args[0].(*StringInstance).Value))
		case 2:
			exception := Send(args[0], "new", NULL, args[1])
			Raise(exception.(object.EmeraldError))
		}

		return NULL
	}
}

func kernelRequireRelative() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		path := args[0].(*StringInstance).Value
		fileParts := strings.Split(ctx.File, "/")
		fileParts[len(fileParts)-1] = ""
		dir := filepath.Join(fileParts...)

		log.InternalDebugF("Attempting to require %s from dir %s", path, dir)

		absoluteFilePath, err := filepath.Abs("/" + filepath.Join(dir, path))
		if err != nil {
			panic(err)
		}

		absolutePathStr := NewString(absoluteFilePath)

		requiredFiles := requiredFilesHash()

		// File has already been loaded
		if requiredFiles.Get(absolutePathStr) != nil {
			log.InternalDebugF("Kernel#require_relative - File %s is already loaded, skipping", absoluteFilePath)
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

		log.InternalDebugF("Kernel#require_relative - Successfully compiled file %s", absoluteFilePath)

		requiredBlock := object.NewClosedBlock(&object.Context{
			Outer: nil,
			File:  absoluteFilePath,
			Self:  MainObject,
			Block: NULL,
			Yield: ctx.Yield,
			BlockGiven: func() bool {
				return false
			},
		}, &object.Block{Instructions: instructions}, []object.EmeraldValue{}, "")

		object.EvalBlock(requiredBlock)

		requiredFiles.Set(absolutePathStr, TRUE)

		return TRUE
	}
}

func kernelInspect() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		return NewString(ctx.Self.Inspect())
	}
}

func requiredFilesHash() *HashInstance {
	requiredFiles := Kernel.InstanceVariableGet("required_files", Kernel, Kernel)
	if requiredFiles == nil {
		hash := NewHash()
		Kernel.InstanceVariableSet("required_files", hash)
		return hash
	} else {
		return requiredFiles.(*HashInstance)
	}
}

func writeToStdout(str string) object.EmeraldError {
	_, err := os.Stdout.WriteString(str)
	if err != nil {
		return raiseStdoutWriteFailed()
	}

	_, err = os.Stdout.WriteString("\n")
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
