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

type IOInstance struct {
	*object.Instance
	FileDescriptor uintptr
	Closed         bool
}

func (rt *Runtime) NewIO(fd uintptr) *IOInstance {
	return &IOInstance{
		Instance:       rt.IO.New(),
		FileDescriptor: fd,
	}
}

func (rt *Runtime) InitIO() {
	rt.IO = rt.DefineClass("IO", rt.Object)

	rt.DefineSingletonMethod(rt.IO, "new", rt.ioNew())
	rt.DefineSingletonMethod(rt.IO, "sysopen", rt.ioSysopen())
	rt.DefineSingletonMethod(rt.IO, "open", rt.ioOpen())
	rt.DefineSingletonMethod(rt.IO, "read", rt.ioRead())

	rt.DefineMethod(rt.IO, "close", rt.ioClose())
	rt.DefineMethod(rt.IO, "getbyte", rt.ioGetbyte())
}

func (rt *Runtime) ioNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		fd := args[0].(*IntegerInstance).Value

		return rt.NewIO(uintptr(fd))
	}
}

func (rt *Runtime) ioSysopen() object.BuiltInMethod {
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
			panic(fmt.Sprintf("rt.IO.sysopen: %s (%q)", err, resolvedPath))
			return rt.Raise(rt.newArgumentError(fmt.Sprintf("%s (%q)", err, resolvedPath)))
		}

		return rt.NewInteger(int64(fd))
	}
}

func (rt *Runtime) ioOpen() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		ioInstance := rt.Send(rt.IO, "new", rt.NULL, kwargs, args...)

		if !ctx.BlockGiven() {
			return ioInstance
		}

		blockResult := ctx.Yield(map[string]object.EmeraldValue{}, ioInstance)

		rt.Send(ioInstance, "close", rt.NULL, map[string]object.EmeraldValue{})

		return blockResult
	}
}

func (rt *Runtime) ioRead() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		fd := rt.Send(rt.IO, "sysopen", rt.NULL, kwargs, args...).(*IntegerInstance).Value

		file := os.NewFile(uintptr(fd), "filename")

		content, err := io.ReadAll(file)
		if err != nil {
			return rt.RaiseGoError(err)
		}

		return rt.NewString(string(content))
	}
}

func (rt *Runtime) ioClose() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		ioInstance := ctx.Self.(*IOInstance)

		if !ioInstance.Closed {
			err := syscall.Close(int(ioInstance.FileDescriptor))
			if err != nil {
				panic(err)
			}
		}

		return rt.NULL
	}
}

func (rt *Runtime) ioGetbyte() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		fd := ctx.Self.(*IOInstance).FileDescriptor

		buffer := make([]byte, 1)

		if n, err := syscall.Read(int(fd), buffer); err != nil {
			panic(err)
		} else if n != 1 {
			panic(fmt.Errorf("expected to read 1 byte but got %d", n))
		}

		return rt.NewInteger(int64(buffer[0]))
	}
}
