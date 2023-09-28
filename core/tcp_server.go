package core

import (
	"context"
	"emerald/object"
	"fmt"
	"net"
	"net/http"
)

var TCPServer *object.Class

type TCPServerInstance struct {
	*object.Instance
	Address  string
	Listener net.Listener
}

func InitTCPServer() {
	TCPServer = DefineClass("TCPServer", Object)

	DefineSingletonMethod(TCPServer, "new", tcpServerNew())

	DefineMethod(TCPServer, "accept", tcpServerAccept())
	DefineMethod(TCPServer, "super_serve", tcpServerSuperServe())
}

func NewTCPServer() *TCPServerInstance {
	return &TCPServerInstance{
		Instance: TCPServer.New(),
		Address:  "",
		Listener: nil,
	}
}

func tcpServerNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := EnforceArity(args, kwargs, 2, 2); err != nil {
			return err
		}

		host, err := EnforceArgumentType[*StringInstance](String, args[0])
		if err != nil {
			return err
		}

		port, err := EnforceArgumentType[*IntegerInstance](Integer, args[1])
		if err != nil {
			return err
		}

		server := NewTCPServer()

		server.Address = fmt.Sprintf("%s:%d", host.Value, port.Value)

		return server
	}
}

func ensureListenerSet(server *TCPServerInstance) object.EmeraldError {
	if server.Listener == nil {
		listener, err := net.Listen("tcp", server.Address)
		if err != nil {
			return RaiseGoError(err)
		}

		server.Listener = listener
	}

	return nil
}

func tcpServerAccept() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		server := ctx.Self.(*TCPServerInstance)

		if err := ensureListenerSet(server); err != nil {
			return err
		}

		conn, err := server.Listener.Accept()
		if err != nil {
			return RaiseGoError(err)
		}

		return NewTCPSocket(conn)
	}
}

type superServer struct {
	ctx *object.Context
}

type contextKey struct {
	key string
}

var ConnContextKey = &contextKey{"http-conn"}

func SaveConnInContext(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, ConnContextKey, c)
}
func GetConn(r *http.Request) net.Conn {
	return r.Context().Value(ConnContextKey).(net.Conn)
}

func (s *superServer) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	s.ctx.Yield(map[string]object.EmeraldValue{}, NewTCPSocket(GetConn(req)))
}

func tcpServerSuperServe() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if !ctx.BlockGiven() {
			return Raise(newArgumentError("You must pass a block to TCPServer#super_serve"))
		}

		emeraldServer := ctx.Self.(*TCPServerInstance)
		if err := ensureListenerSet(emeraldServer); err != nil {
			return err
		}

		goServer := &http.Server{
			Handler:     &superServer{ctx: ctx},
			ConnContext: SaveConnInContext,
		}

		err := goServer.Serve(emeraldServer.Listener)
		if err != nil {
			return RaiseGoError(err)
		}

		return NULL
	}
}
