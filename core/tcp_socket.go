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

func (rt *Runtime) NewTCPSocket(conn net.Conn) *TCPSocketInstance {
	return &TCPSocketInstance{
		Instance: rt.TCPSocket.New(),
		Conn:     conn,
		tp:       textproto.NewReader(bufio.NewReader(conn)),
		timeout:  &tcpSocketTimeout{},
	}
}

func (rt *Runtime) tcpSocketTimeoutGet() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		socket := ctx.Self.(*TCPSocketInstance)

		if !socket.timeout.IsValid {
			return rt.NULL
		}

		return rt.NewInteger(int64(socket.timeout.Get() / time.Millisecond))
	}
}

func (rt *Runtime) tcpSocketTimeoutSet() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		rt.EnforceArity(args, kwargs, 1, 1)

		socket := ctx.Self.(*TCPSocketInstance)

		if args[0] == rt.NULL {
			socket.timeout.Reset()
			return rt.NULL
		}

		newValue, err := EnforceArgumentType[*IntegerInstance](rt, rt.Integer, args[0])
		if err != nil {
			return err
		}

		socket.timeout.Set(time.Duration(newValue.Value) * time.Millisecond)

		return newValue
	}
}

func (rt *Runtime) tcpSocketGets() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		socket := ctx.Self.(*TCPSocketInstance)

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

		socket := ctx.Self.(*TCPSocketInstance)

		content, emeraldErr := EnforceArgumentType[*StringInstance](rt, rt.String, args[0])
		if emeraldErr != nil {
			return emeraldErr
		}

		debug.InternalDebugF("Writing to TCPSocket: %s", content.Value)

		bytesWritten, err := socket.Conn.Write([]byte(content.Value))
		if err != nil {
			return rt.RaiseGoError(err)
		}

		return rt.NewInteger(int64(bytesWritten))
	}
}

func (rt *Runtime) tcpSocketClose() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		socket := ctx.Self.(*TCPSocketInstance)

		err := socket.Conn.Close()
		if err != nil {
			return rt.RaiseGoError(err)
		}

		return rt.NULL
	}
}
