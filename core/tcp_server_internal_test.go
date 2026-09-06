package core

import "testing"

func TestTCPServerRejectsUninitializedListener(t *testing.T) {
	rt := NewRuntime()
	rt.Init()
	server := &TCPServerInstance{}
	err := rt.ensureListenerSet(server)
	if server.Listener != nil {
		server.Listener.Close()
		t.Fatal("uninitialized TCPServer opened a listener")
	}
	if err == nil || err.Message() != "uninitialized TCPServer" || !rt.ExceptionIsRaised() {
		t.Fatalf("expected an Emerald exception before listening, got %v", err)
	}
}
