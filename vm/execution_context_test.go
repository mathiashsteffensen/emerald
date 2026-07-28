package vm

import (
	"context"
	"emerald/bytecode"
	"emerald/compiler"
	"emerald/core"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestRunContextStopsExecutionWithoutRaising(t *testing.T) {
	tests := map[string]string{
		"bytecode loop": "while true\nend",
		"native loop":   "9223372036854775807.times { |i| i }",
	}

	for name, source := range tests {
		t.Run(name, func(t *testing.T) {
			machine, rt := compileContextTestVM(source)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
			defer cancel()

			err := machine.RunContext(ctx)
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("expected deadline exceeded, got %v", err)
			}
			if rt.ExceptionIsRaised() {
				t.Fatal("execution cancellation must not raise an Emerald exception")
			}
		})
	}
}

func TestRunContextHonorsCancellationBeforeFirstInstruction(t *testing.T) {
	machine, rt := compileContextTestVM("$started = true")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := machine.RunContext(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
	if started := rt.Heap.GetGlobalVariableString("$started"); started.IsDefined() {
		t.Fatalf("expected no instructions to run, got %s", started.Inspect())
	}
}

func TestRunContextBoundsRegexpMatching(t *testing.T) {
	source := fmt.Sprintf("%q =~ /(a+)+$/", strings.Repeat("a", 2_000)+"!")
	machine, rt := compileContextTestVM(source)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := machine.RunContext(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
	if rt.ExceptionIsRaised() {
		t.Fatal("regexp cancellation must not raise an Emerald exception")
	}
}

func compileContextTestVM(source string) (*VM, *core.Runtime) {
	rt := core.NewRuntime()
	rt.Init()
	rt.CompileBlock = func(fileName string, content string) *bytecode.Bytecode {
		return compiler.CompileBlock(fileName, content, rt)
	}

	bc := compiler.Compile("test", source, rt)
	return New("test", bc, rt), rt
}
