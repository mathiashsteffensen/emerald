package core

import (
	"emerald/object"
	"fmt"
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

		fd, err := syscall.Open(path, syscall.O_NONBLOCK, 0)
		if err != nil {
			panic(err)
		}

		return NewInteger(int64(fd))
	}
}

func ioOpen() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		io := Send(IO, "new", NULL, args...)

		if !ctx.BlockGiven() {
			return io
		}

		blockResult := ctx.Yield(io)

		Send(io, "close", NULL)

		return blockResult
	}
}

func ioClose() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		io := ctx.Self.(*IOInstance)

		if !io.Closed {
			err := syscall.Close(int(io.FileDescriptor))
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
