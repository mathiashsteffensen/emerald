package core

import (
	"emerald/object"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

var IO *object.Class

type IOInstance struct {
	*object.Instance
	FileDescriptor uintptr
	Closed         bool
}

func NewIO(fd uintptr) *IOInstance {
	return &IOInstance{
		Instance:       IO.New(),
		FileDescriptor: fd,
	}
}

func InitIO() {
	IO = DefineClass("IO", Object)

	DefineSingletonMethod(IO, "new", ioNew())
	DefineSingletonMethod(IO, "sysopen", ioSysopen())
	DefineSingletonMethod(IO, "open", ioOpen())
	DefineSingletonMethod(IO, "read", ioRead())

	DefineMethod(IO, "close", ioClose())
	DefineMethod(IO, "getbyte", ioGetbyte())
}

func ioNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		fd := args[0].(*IntegerInstance).Value

		return NewIO(uintptr(fd))
	}
}

func ioSysopen() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		path := args[0].(*StringInstance).Value

		var resolvedPath string
		if filepath.IsAbs(path) {
			resolvedPath = path
		} else {
			fileParts := strings.Split(ctx.File, "/")
			fileParts[len(fileParts)-1] = path
			resolvedPath = filepath.Join(fileParts...)
		}

		fd, err := syscall.Open("/"+resolvedPath, syscall.O_NONBLOCK, 0)
		if err != nil {
			panic(fmt.Sprintf("IO.sysopen: %s (%q)", err, resolvedPath))
			return Raise(newArgumentError(fmt.Sprintf("%s (%q)", err, resolvedPath)))
		}

		return NewInteger(int64(fd))
	}
}

func ioOpen() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		ioInstance := Send(IO, "new", NULL, kwargs, args...)

		if !ctx.BlockGiven() {
			return ioInstance
		}

		blockResult := ctx.Yield(map[string]object.EmeraldValue{}, ioInstance)

		Send(ioInstance, "close", NULL, map[string]object.EmeraldValue{})

		return blockResult
	}
}

func ioRead() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		fd := Send(IO, "sysopen", NULL, kwargs, args...).(*IntegerInstance).Value

		file := os.NewFile(uintptr(fd), "filename")

		content, err := io.ReadAll(file)
		if err != nil {
			return RaiseGoError(err)
		}

		return NewString(string(content))
	}
}

func ioClose() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		ioInstance := ctx.Self.(*IOInstance)

		if !ioInstance.Closed {
			err := syscall.Close(int(ioInstance.FileDescriptor))
			if err != nil {
				panic(err)
			}
		}

		return NULL
	}
}

func ioGetbyte() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		fd := ctx.Self.(*IOInstance).FileDescriptor

		buffer := make([]byte, 1)

		if n, err := syscall.Read(int(fd), buffer); err != nil {
			panic(err)
		} else if n != 1 {
			panic(fmt.Errorf("expected to read 1 byte but got %d", n))
		}

		return NewInteger(int64(buffer[0]))
	}
}
