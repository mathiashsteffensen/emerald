package compiler

import (
	"emerald/core"
	"testing"
)

func TestCompile(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Compiling panicked %s", r)
		}
	}()

	rt := core.NewRuntime()
	rt.Init()
	CompileToBytecode("test.rb", "puts(\"Hello\")", rt)
}
