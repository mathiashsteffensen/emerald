package core

import (
	"bufio"
	"emerald/debug"
	"emerald/object"
	"net"
	"net/textproto"
	"time"
)

type tcpSocketTimeout struct {
	time.Duration
	IsValid bool
}

func (timeout *tcpSocketTimeout) Set(value time.Duration) {
	timeout.Duration = value
	timeout.IsValid = true
}

func (timeout *tcpSocketTimeout) Get() time.Duration {
	return timeout.Duration
}

func (timeout *tcpSocketTimeout) Reset() {
	timeout.IsValid = false
	timeout.Duration = time.Duration(0)
}

type TCPSocketInstance struct {
	*object.Instance
	net.Conn
	tp      *textproto.Reader
	timeout *tcpSocketTimeout
}

func (rt *Runtime) InitTCPSocket() {
	rt.TCPSocket = rt.DefineClass("TCPSocket", rt.Object)

	rt.DefineMethod(rt.TCPSocket, "gets", rt.tcpSocketGets())
	rt.DefineMethod(rt.TCPSocket, "write", rt.tcpSocketWrite())
	rt.DefineMethod(rt.TCPSocket, "close", rt.tcpSocketClose())
	rt.DefineMethod(rt.TCPSocket, "timeout", rt.tcpSocketTimeoutGet())
	rt.DefineMethod(rt.TCPSocket, "timeout=", rt.tcpSocketTimeoutSet())
}

func (rt *Runtime) NewTCPSocket(conn net.Conn) object.EmeraldValue {
	return object.NewHeapObject(&TCPSocketInstance{
		Instance: rt.TCPSocket.Heap.(*object.Class).New(),
		Conn:     conn,
		tp:       textproto.NewReader(bufio.NewReader(conn)),
		timeout:  &tcpSocketTimeout{},
	})
}

func (rt *Runtime) tcpSocketTimeoutGet() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		socket := ctx.Self.Heap.(*TCPSocketInstance)

		if !socket.timeout.IsValid {
			return rt.NULL
		}

		return rt.NewInteger(int64(socket.timeout.Get() / time.Millisecond))
	}
}

func (rt *Runtime) tcpSocketTimeoutSet() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		rt.EnforceArity(args, kwargs, 1, 1)

		socket := ctx.Self.Heap.(*TCPSocketInstance)

		if args[0].IsNil() {
			socket.timeout.Reset()
			return rt.NULL
		}

		newValue, err := rt.EnforceIntegerArg(args[0])
		if err != nil {
			return object.NewHeapObject(err)
		}

		socket.timeout.Set(time.Duration(newValue) * time.Millisecond)

		return args[0]
	}
}

func (rt *Runtime) tcpSocketGets() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		socket := ctx.Self.Heap.(*TCPSocketInstance)

		line, err := socket.tp.ReadLine()
		if err != nil {
			debug.DebugF("Error reading from socket: %s", err)
			return rt.NULL
		}

		if line == "" {
			return rt.NULL
		}

		return rt.NewString(line)
	}
}

func (rt *Runtime) tcpSocketWrite() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		rt.EnforceArity(args, kwargs, 1, 1)

		socket := ctx.Self.Heap.(*TCPSocketInstance)

		content, emeraldErr := EnforceArgumentType[*StringInstance](rt, rt.String, args[0])
		if emeraldErr != nil {
			return object.NewHeapObject(emeraldErr)
		}

		debug.InternalDebugF("Writing to TCPSocket: %s", content.Value)

		bytesWritten, err := socket.Conn.Write([]byte(content.Value))
		if err != nil {
			return object.NewHeapObject(rt.RaiseGoError(err))
		}

		return rt.NewInteger(int64(bytesWritten))
	}
}

func (rt *Runtime) tcpSocketClose() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		socket := ctx.Self.Heap.(*TCPSocketInstance)

		err := socket.Conn.Close()
		if err != nil {
			return object.NewHeapObject(rt.RaiseGoError(err))
		}

		return rt.NULL
	}
}
