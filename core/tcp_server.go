package core

import (
	"context"
	"emerald/object"
	"fmt"
	"net"
	"net/http"
)

type TCPServerInstance struct {
	*object.Instance
	Address  string
	Listener net.Listener
}

func (rt *Runtime) InitTCPServer() {
	rt.TCPServer = rt.DefineClass("TCPServer", rt.Object)

	rt.DefineSingletonMethod(rt.TCPServer, "new", rt.tcpServerNew())

	rt.DefineMethod(rt.TCPServer, "accept", rt.tcpServerAccept())
	rt.DefineMethod(rt.TCPServer, "super_serve", rt.tcpServerSuperServe())
}

func (rt *Runtime) NewTCPServer() object.EmeraldValue {
	return object.NewHeapObject(&TCPServerInstance{
		Instance: rt.TCPServer.Heap.(*object.Class).New(),
		Address:  "",
		Listener: nil,
	})
}

func (rt *Runtime) tcpServerNew() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if _, err := rt.EnforceArity(args, kwargs, 2, 2); err != nil {
			return object.NewHeapObject(err)
		}

		host, err := EnforceArgumentType[*StringInstance](rt, rt.String, args[0])
		if err != nil {
			return object.NewHeapObject(err)
		}

		port, err := rt.EnforceIntegerArg(args[1])
		if err != nil {
			return object.NewHeapObject(err)
		}

		serverVal := rt.NewTCPServer()
		server := serverVal.Heap.(*TCPServerInstance)

		server.Address = fmt.Sprintf("%s:%d", host.Value, port)

		return serverVal
	}
}

func (rt *Runtime) ensureListenerSet(server *TCPServerInstance) object.EmeraldError {
	if server.Listener == nil {
		listener, err := net.Listen("tcp", server.Address)
		if err != nil {
			return rt.RaiseGoError(err)
		}

		server.Listener = listener
	}

	return nil
}

func (rt *Runtime) tcpServerAccept() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		server := ctx.Self.Heap.(*TCPServerInstance)

		if err := rt.ensureListenerSet(server); err != nil {
			return object.NewHeapObject(err)
		}

		conn, err := server.Listener.Accept()
		if err != nil {
			return object.NewHeapObject(rt.RaiseGoError(err))
		}

		return rt.NewTCPSocket(conn)
	}
}

type superServer struct {
	rt  *Runtime
	ctx *object.Context
}

type contextKey struct {
	key string
}

var ConnContextKey = &contextKey{"http-conn"}

func (rt *Runtime) SaveConnInContext(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, ConnContextKey, c)
}
func (rt *Runtime) GetConn(r *http.Request) net.Conn {
	return r.Context().Value(ConnContextKey).(net.Conn)
}

func (s *superServer) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	s.ctx.Yield(map[string]object.EmeraldValue{}, s.rt.NewTCPSocket(s.rt.GetConn(req)))
}

func (rt *Runtime) tcpServerSuperServe() object.BuiltInMethod {
	return func(ctx *object.Context, kwargs map[string]object.EmeraldValue, args ...object.EmeraldValue) object.EmeraldValue {
		if !ctx.BlockGiven() {
			return object.NewHeapObject(rt.Raise(rt.newArgumentError("You must pass a block to rt.TCPServer#super_serve")))
		}

		emeraldServer := ctx.Self.Heap.(*TCPServerInstance)
		if err := rt.ensureListenerSet(emeraldServer); err != nil {
			return object.NewHeapObject(err)
		}

		goServer := &http.Server{
			Handler:     &superServer{ctx: ctx, rt: rt},
			ConnContext: rt.SaveConnInContext,
		}

		err := goServer.Serve(emeraldServer.Listener)
		if err != nil {
			return object.NewHeapObject(rt.RaiseGoError(err))
		}

		return rt.NULL
	}
}
