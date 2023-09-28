package core

import (
	"bufio"
	"emerald/debug"
	"emerald/object"
	"net"
	"net/textproto"
	"time"
)

var TCPSocket *object.Class

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

func InitTCPSocket() {
	TCPSocket = DefineClass("TCPSocket", Object)

	DefineMethod(TCPSocket, "gets", tcpSocketGets())
	DefineMethod(TCPSocket, "write", tcpSocketWrite())
	DefineMethod(TCPSocket, "close", tcpSocketClose())
	DefineMethod(TCPSocket, "timeout", tcpSocketTimeoutGet())
	DefineMethod(TCPSocket, "timeout=", tcpSocketTimeoutSet())
}

func NewTCPSocket(conn net.Conn) *TCPSocketInstance {
	return &TCPSocketInstance{
		Instance: TCPSocket.New(),
		Conn:     conn,
		tp:       textproto.NewReader(bufio.NewReader(conn)),
		timeout:  &tcpSocketTimeout{},
	}
}

func tcpSocketTimeoutGet() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		socket := ctx.Self.(*TCPSocketInstance)

		if !socket.timeout.IsValid {
			return NULL
		}

		return NewInteger(int64(socket.timeout.Get() / time.Millisecond))
	}
}

func tcpSocketTimeoutSet() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		EnforceArity(args, kwargs, 1, 1)

		socket := ctx.Self.(*TCPSocketInstance)

		if args[0] == NULL {
			socket.timeout.Reset()
			return NULL
		}

		newValue, err := EnforceArgumentType[*IntegerInstance](Integer, args[0])
		if err != nil {
			return err
		}

		socket.timeout.Set(time.Duration(newValue.Value) * time.Millisecond)

		return newValue
	}
}

func tcpSocketGets() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		socket := ctx.Self.(*TCPSocketInstance)

		line, err := socket.tp.ReadLine()
		if err != nil {
			debug.DebugF("Error reading from socket: %s", err)
			return NULL
		}

		if line == "" {
			return NULL
		}

		return NewString(line)
	}
}

func tcpSocketWrite() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		EnforceArity(args, kwargs, 1, 1)

		socket := ctx.Self.(*TCPSocketInstance)

		content, emeraldErr := EnforceArgumentType[*StringInstance](String, args[0])
		if emeraldErr != nil {
			return emeraldErr
		}

		debug.InternalDebugF("Writing to TCPSocket: %s", content.Value)

		bytesWritten, err := socket.Conn.Write([]byte(content.Value))
		if err != nil {
			return RaiseGoError(err)
		}

		return NewInteger(int64(bytesWritten))
	}
}

func tcpSocketClose() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		socket := ctx.Self.(*TCPSocketInstance)

		err := socket.Conn.Close()
		if err != nil {
			return RaiseGoError(err)
		}

		return NULL
	}
}
