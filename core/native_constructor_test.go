package core_test

import (
	"emerald/compiler"
	"emerald/core"
	"emerald/object"
	"emerald/vm"
	"fmt"
	"os"
	"testing"
)

func TestInheritedNativeResourceStorage(t *testing.T) {
	t.Run("IO is closed without consuming or closing descriptors", func(t *testing.T) {
		reader, writer, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		defer reader.Close()
		defer writer.Close()
		rt := core.NewRuntime()
		rt.Init()
		code := fmt.Sprintf(`class Custom < IO; def initialize(fd); end; end; Custom.new(%d)`, reader.Fd())
		machine := vm.New("test", compiler.Compile("test", code, rt), rt)
		machine.Run()
		value := machine.LastPoppedStackElem()
		io, ok := value.Heap.(*core.IOInstance)
		if !ok || !io.Closed || io.FileDescriptor != ^uintptr(0) {
			t.Fatalf("uninitialized IO must be closed with an invalid descriptor: %#v", value.Heap)
		}
		if !rt.Send(value, "close", rt.NULL, nil).IsNil() || rt.ExceptionIsRaised() {
			t.Fatal("closing uninitialized IO must be harmless")
		}
		if _, err := writer.Write([]byte{42}); err != nil {
			t.Fatal(err)
		}
		rt.Send(value, "getbyte", rt.NULL, nil)
		raised := rt.Heap.GetGlobalVariableString("$!")
		if err, ok := raised.Heap.(object.EmeraldError); !ok || err.Message() != "closed stream" {
			t.Fatalf("getbyte must raise an Emerald closed-stream exception: %v", raised)
		}
		data := make([]byte, 1)
		if n, err := reader.Read(data); err != nil || n != 1 || data[0] != 42 {
			t.Fatalf("descriptor was consumed or closed: %v, %v", data, err)
		}
	})
	t.Run("TCPServer cannot listen without initialization", func(t *testing.T) {
		rt := core.NewRuntime()
		rt.Init()
		machine := vm.New("test", compiler.Compile("test", `class Custom < TCPServer; end; Custom.new`, rt), rt)
		machine.Run()
		if rt.ExceptionIsRaised() {
			t.Fatalf("subclass allocation failed: %s", rt.Heap.GetGlobalVariableString("$!").Inspect())
		}
		value := machine.LastPoppedStackElem()
		server, ok := value.Heap.(*core.TCPServerInstance)
		if !ok || server.Address != "" || server.Listener != nil {
			t.Fatalf("unexpected uninitialized server: %#v", value.Heap)
		}
	})
}

func TestNativeConstructorDirectOwners(t *testing.T) {
	for _, name := range []string{"Exception", "StandardError", "ArgumentError", "TypeError", "NameError", "NoMethodError", "RuntimeError", "LoadError"} {
		runCoreTests(t, []coreTestCase{{name: name, input: name + `.new("boom").to_s`, expected: "boom"}})
	}
	runCoreTests(t, []coreTestCase{
		{name: "regexp", input: `Regexp.new("abc") === "abc"`, expected: true},
		{name: "range", input: `values = []; Range.new(1, 2).each { |n| values << n }; values`, expected: []any{1, 2}},
		{name: "string", input: `String.new("abc")`, expected: "abc"},
		{name: "time", input: `Time.new.class == Time`, expected: true},
		{name: "time now", input: `Time.now.class == Time`, expected: true},
		{name: "tcp server", input: `TCPServer.new("127.0.0.1", 0).class == TCPServer`, expected: true},
		{name: "class", input: `Class.new.class == Class`, expected: true},
	})
}
