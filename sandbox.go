package emerald

import (
	"context"
	"emerald/object"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	ErrSandboxTimeout           = errors.New("emerald: sandbox timeout")
	ErrSandboxWorkerFailed      = errors.New("emerald: sandbox worker failed")
	ErrSandboxUnsupportedResult = errors.New("emerald: sandbox result is not transferable")
)

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

// SandboxWorkerError reports a worker process or protocol failure. A worker
// can terminate because of a VM defect or a resource limit, neither of which
// is reliably distinguishable on every supported platform.
type SandboxWorkerError struct {
	Phase  string
	Detail string
	Cause  error
}

func (e *SandboxWorkerError) Error() string {
	message := ErrSandboxWorkerFailed.Error()
	if e.Phase != "" {
		message += " during " + e.Phase
	}
	if e.Detail != "" {
		message += ": " + e.Detail
	}
	return message
}

func (e *SandboxWorkerError) Unwrap() []error {
	if e.Cause == nil {
		return []error{ErrSandboxWorkerFailed}
	}
	return []error{ErrSandboxWorkerFailed, e.Cause}
}

// SandboxUnsupportedResultError identifies a value that cannot cross the
// worker process boundary.
type SandboxUnsupportedResultError struct {
	Type   string
	Reason string
}

func (e *SandboxUnsupportedResultError) Error() string {
	message := ErrSandboxUnsupportedResult.Error()
	if e.Type != "" {
		message += ": " + e.Type
	}
	if e.Reason != "" {
		message += " (" + e.Reason + ")"
	}
	return message
}

func (e *SandboxUnsupportedResultError) Unwrap() error {
	return ErrSandboxUnsupportedResult
}

type SandboxOptions struct {
	Timeout          time.Duration
	MemoryLimitBytes int64
	WorkerPath       string
}

// Sandbox is an immutable execution policy. Each evaluation uses a fresh,
// restricted worker process so successful and interrupted scripts cannot
// share state.
type Sandbox struct {
	timeout          time.Duration
	memoryLimitBytes int64
	workerPath       string
}

func NewSandbox(options SandboxOptions) (*Sandbox, error) {
	if options.Timeout <= 0 {
		return nil, fmt.Errorf("emerald: sandbox timeout must be positive")
	}
	if options.MemoryLimitBytes <= 0 {
		return nil, fmt.Errorf("emerald: sandbox memory limit must be positive")
	}
	if options.WorkerPath == "" {
		return nil, fmt.Errorf("emerald: sandbox worker path is required")
	}

	workerPath, err := exec.LookPath(options.WorkerPath)
	if err != nil {
		return nil, fmt.Errorf("emerald: sandbox worker: %w", err)
	}
	workerPath, err = filepath.Abs(workerPath)
	if err != nil {
		return nil, fmt.Errorf("emerald: sandbox worker path: %w", err)
	}
	info, err := os.Stat(workerPath)
	if err != nil {
		return nil, fmt.Errorf("emerald: sandbox worker: %w", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("emerald: sandbox worker path is a directory")
	}

	return &Sandbox{
		timeout:          options.Timeout,
		memoryLimitBytes: options.MemoryLimitBytes,
		workerPath:       workerPath,
	}, nil
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

	value, phase, err := s.evalWorker(ctx, fileName, content)
	if ctxErr := ctx.Err(); ctxErr != nil {
		if errors.Is(ctxErr, context.DeadlineExceeded) && sandboxOwnsDeadline {
			return object.EmeraldValue{}, &SandboxTimeoutError{Timeout: s.timeout, Phase: phase}
		}
		return object.EmeraldValue{}, ctxErr
	}
	if errors.Is(err, context.DeadlineExceeded) {
		if !sandboxOwnsDeadline {
			return object.EmeraldValue{}, context.DeadlineExceeded
		}
		return object.EmeraldValue{}, &SandboxTimeoutError{Timeout: s.timeout, Phase: phase}
	}

	return value, err
}
