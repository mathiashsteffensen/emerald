package emerald

import (
	"context"
	"emerald/core"
	"emerald/object"
	"errors"
	"fmt"
	"time"
)

var ErrSandboxTimeout = errors.New("emerald: sandbox timeout")

type SandboxTimeoutError struct {
	Timeout time.Duration
	Phase   string
}

func (e *SandboxTimeoutError) Error() string {
	if e.Phase != "" {
		return fmt.Sprintf("%s during %s after %s", ErrSandboxTimeout, e.Phase, e.Timeout)
	}
	return fmt.Sprintf("%s after %s", ErrSandboxTimeout, e.Timeout)
}

func (e *SandboxTimeoutError) Unwrap() error {
	return ErrSandboxTimeout
}

func (e *SandboxTimeoutError) Is(target error) bool {
	return target == context.DeadlineExceeded
}

type SandboxOptions struct {
	Timeout time.Duration
}

// Sandbox is an immutable execution policy. Each evaluation uses a fresh,
// restricted runtime so successful and interrupted scripts cannot share state.
type Sandbox struct {
	timeout time.Duration
}

func NewSandbox(options SandboxOptions) (*Sandbox, error) {
	if options.Timeout <= 0 {
		return nil, fmt.Errorf("emerald: sandbox timeout must be positive")
	}

	return &Sandbox{timeout: options.Timeout}, nil
}

func (s *Sandbox) Eval(content string) (object.EmeraldValue, error) {
	return s.EvalContext(context.Background(), content)
}

func (s *Sandbox) EvalContext(ctx context.Context, content string) (object.EmeraldValue, error) {
	return s.evalFileContext(ctx, "(sandbox)", content)
}

func (s *Sandbox) EvalFile(fileName string, content string) (object.EmeraldValue, error) {
	return s.EvalFileContext(context.Background(), fileName, content)
}

func (s *Sandbox) EvalFileContext(
	parent context.Context,
	fileName string,
	content string,
) (object.EmeraldValue, error) {
	return s.evalFileContext(parent, fileName, content)
}

func (s *Sandbox) evalFileContext(
	parent context.Context,
	fileName string,
	content string,
) (object.EmeraldValue, error) {
	if parent == nil {
		return object.EmeraldValue{}, fmt.Errorf("emerald: nil sandbox context")
	}
	if err := parent.Err(); err != nil {
		return object.EmeraldValue{}, err
	}

	parentDeadline, parentHasDeadline := parent.Deadline()

	ctx, cancel := context.WithTimeout(parent, s.timeout)
	defer cancel()
	effectiveDeadline, _ := ctx.Deadline()
	sandboxOwnsDeadline := !parentHasDeadline || effectiveDeadline.Before(parentDeadline)

	rt := core.NewRuntime()
	rt.InitSandbox()

	value, err := newEngine(rt).evalFileContext(ctx, fileName, content, EvalOptions{})
	if errors.Is(err, context.DeadlineExceeded) {
		if !sandboxOwnsDeadline {
			return object.EmeraldValue{}, context.DeadlineExceeded
		}

		var phaseErr *evaluationPhaseError
		errors.As(err, &phaseErr)

		timeoutErr := &SandboxTimeoutError{Timeout: s.timeout}
		if phaseErr != nil {
			timeoutErr.Phase = phaseErr.phase
		}
		return object.EmeraldValue{}, timeoutErr
	}

	return value, err
}
