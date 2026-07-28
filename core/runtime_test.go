package core_test

import (
	"emerald/core"
	"testing"
)

var hostConstantNames = []string{
	"IO",
	"Dir",
	"File",
	"Time",
	"TCPServer",
	"TCPSocket",
}

var hostKernelMethodNames = []string{
	"require_relative",
	"sleep",
	"puts",
	"print",
}

func TestRuntimeInitSandboxExcludesHostCapabilities(t *testing.T) {
	rt := core.NewRuntime()
	rt.InitSandbox()

	for _, name := range hostConstantNames {
		if value := rt.Object.NamespaceDefinitionGet(name); value.IsDefined() {
			t.Errorf("expected %s to be undefined in sandbox runtime", name)
		}
	}

	for _, name := range hostKernelMethodNames {
		if _, ok := rt.Kernel.BuiltInMethodSet()[name]; ok {
			t.Errorf("expected Kernel#%s to be undefined in sandbox runtime", name)
		}
		if _, ok := rt.Kernel.SingletonClass().BuiltInMethodSet()[name]; ok {
			t.Errorf("expected Kernel.%s to be undefined in sandbox runtime", name)
		}
	}

	if _, ok := rt.Kernel.BuiltInMethodSet()["raise"]; !ok {
		t.Error("expected pure Kernel#raise to remain defined in sandbox runtime")
	}

	for _, name := range []string{"Object", "String", "Array", "Hash", "Regexp"} {
		if value := rt.Object.NamespaceDefinitionGet(name); !value.IsDefined() {
			t.Errorf("expected pure core constant %s to be defined in sandbox runtime", name)
		}
	}
}

func TestRuntimeInitIncludesHostCapabilities(t *testing.T) {
	rt := core.NewRuntime()
	rt.Init()

	for _, name := range hostConstantNames {
		if value := rt.Object.NamespaceDefinitionGet(name); !value.IsDefined() {
			t.Errorf("expected %s to be defined in full runtime", name)
		}
	}

	for _, name := range hostKernelMethodNames {
		if _, ok := rt.Kernel.BuiltInMethodSet()[name]; !ok {
			t.Errorf("expected Kernel#%s to be defined in full runtime", name)
		}
		if _, ok := rt.Kernel.SingletonClass().BuiltInMethodSet()[name]; !ok {
			t.Errorf("expected Kernel.%s to be defined in full runtime", name)
		}
	}
}
