package emerald

import (
	"context"
	"emerald/bytecode"
	"emerald/compiler"
	"emerald/core"
	"emerald/object"
	"emerald/vm"
	"fmt"
)

type Engine struct {
	Runtime *core.Runtime
}

func New() *Engine {
	rt := core.NewRuntime()
	rt.Init()

	return newEngine(rt)
}

func newEngine(rt *core.Runtime) *Engine {
	rt.CompileBlock = func(fileName string, content string) *bytecode.Bytecode {
		return compiler.CompileBlock(fileName, content, rt)
	}

	return &Engine{
		Runtime: rt,
	}
}

// EvalOptions contains host-provided inputs for one evaluation.
// Nil slices are treated as empty; Emerald never reads process arguments
// or load paths implicitly.
type EvalOptions struct {
	Args     []string
	LoadPath []string
}

func (e *Engine) Eval(content string) (object.EmeraldValue, error) {
	return e.EvalWithOptions(content, EvalOptions{})
}

func (e *Engine) EvalWithOptions(content string, options EvalOptions) (object.EmeraldValue, error) {
	return e.EvalFileWithOptions("(irb)", content, options)
}

type EvalError struct {
	ClassName  string
	Message    string
	Stacktrace []byte
}

func (e EvalError) Error() string {
	return fmt.Sprintf("%s: %s", e.ClassName, e.Message)
}

type evaluationPhaseError struct {
	phase string
	err   error
}

func (e *evaluationPhaseError) Error() string {
	return e.err.Error()
}

func (e *evaluationPhaseError) Unwrap() error {
	return e.err
}

func (e *Engine) EvalFile(fileName string, content string) (object.EmeraldValue, error) {
	return e.EvalFileWithOptions(fileName, content, EvalOptions{})
}

func (e *Engine) EvalFileWithOptions(fileName string, content string, options EvalOptions) (object.EmeraldValue, error) {
	return e.evalFileContext(context.Background(), fileName, content, options)
}

func (e *Engine) evalFileContext(
	ctx context.Context,
	fileName string,
	content string,
	options EvalOptions,
) (object.EmeraldValue, error) {
	e.Runtime.Heap.SetGlobalVariableString("$!", object.EmeraldValue{})
	e.Runtime.OnRaise = nil

	bc, err := compiler.CompileContext(ctx, fileName, content, e.Runtime)
	if err != nil {
		return object.EmeraldValue{}, &evaluationPhaseError{phase: "compile", err: err}
	}

	if err := e.evalError(); err != nil {
		return object.EmeraldValue{}, err
	}

	machine := vm.NewWithOptions(fileName, bc, e.Runtime, vm.Options{
		Args:     append([]string(nil), options.Args...),
		LoadPath: append([]string(nil), options.LoadPath...),
	})
	if err := machine.RunContext(ctx); err != nil {
		return object.EmeraldValue{}, &evaluationPhaseError{phase: "execute", err: err}
	}

	if err := e.evalError(); err != nil {
		return object.EmeraldValue{}, err
	}

	return machine.LastPoppedStackElem(), nil
}

func (e *Engine) evalError() error {
	globalException := e.Runtime.Heap.GetGlobalVariableString("$!")
	if !globalException.IsNil() && globalException != e.Runtime.NULL {
		if emErr, ok := globalException.Heap.(object.EmeraldError); ok {
			return EvalError{
				ClassName:  emErr.ClassName(),
				Message:    emErr.Message(),
				Stacktrace: []byte{},
			}
		}

		return fmt.Errorf("non-emerald-error exception: %s", globalException.Inspect())
	}

	return nil
}
