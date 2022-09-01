package core

import (
	"emerald/log"
	"emerald/object"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var Kernel *object.Module

var Compile func(fileName string, content string) []byte

func InitKernel() {
	Kernel = DefineModule(nil, "Kernel")

	DefineMethod(Kernel, "inspect", kernelInspect())
	DefineMethod(Kernel, "require_relative", kernelRequireRelative())
	DefineMethod(Kernel, "class", kernelClass())
	DefineMethod(Kernel, "kind_of?", kernelKindOf())
	DefineMethod(Kernel, "is_a?", kernelKindOf())
	DefineMethod(Kernel, "include", kernelInclude())

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
			return NewTypeError("class or module required")
		}

		for _, ancestor := range ctx.Self.Class().Ancestors() {
			if ancestor == args[0] {
				return TRUE
			}
		}

		return FALSE
	}
}

func kernelPuts() object.BuiltInMethod {
	return func(ctx *object.Context, args ...object.EmeraldValue) object.EmeraldValue {
		var strArr []any

		for _, arg := range args {
			val := Send(arg, "to_s", NULL)

			strArr = append(strArr, val.Inspect())
		}

		fmt.Println(strArr...)

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
				panic(fmt.Errorf("cannot load such file -- %s", absoluteFilePath))
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
		}, &object.Block{Instructions: instructions}, []object.EmeraldValue{})

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
