package emerald

import (
	"context"
	"emerald/core"
	"emerald/internal/sandboxwire"
	"emerald/object"
	"errors"
	"fmt"
	"io"
	"os/exec"
	runtimedebug "runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	maxSandboxProtocolFrameBytes = 16 << 20
	maxSandboxWorkerStderrBytes  = 8 << 10
)

func sandboxFrameLimit(memoryLimitBytes int64) int {
	if memoryLimitBytes <= 0 || memoryLimitBytes >= maxSandboxProtocolFrameBytes {
		return maxSandboxProtocolFrameBytes
	}
	return int(memoryLimitBytes)
}

func (s *Sandbox) evalWorker(
	ctx context.Context,
	fileName string,
	content string,
) (object.EmeraldValue, string, error) {
	cmd := exec.CommandContext(ctx, s.workerPath)
	cmd.Env = []string{
		"GOMEMLIMIT=" + strconv.FormatInt(s.memoryLimitBytes, 10) + "B",
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return object.EmeraldValue{}, "", &SandboxWorkerError{Detail: "open worker stdin", Cause: err}
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return object.EmeraldValue{}, "", &SandboxWorkerError{Detail: "open worker stdout", Cause: err}
	}
	stderr := &sandboxLimitedBuffer{limit: maxSandboxWorkerStderrBytes}
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return object.EmeraldValue{}, "", &SandboxWorkerError{Detail: "start worker", Cause: err}
	}
	if err := applySandboxHardMemoryLimit(cmd.Process.Pid, s.memoryLimitBytes); err != nil {
		stopSandboxWorker(cmd)
		return object.EmeraldValue{}, "", &SandboxWorkerError{Detail: "apply hard memory limit", Cause: err}
	}

	deadline, _ := ctx.Deadline()
	request := sandboxwire.Request{
		Version:          sandboxwire.Version,
		FileName:         fileName,
		Content:          content,
		DeadlineUnixNano: deadline.UnixNano(),
		MemoryLimitBytes: s.memoryLimitBytes,
	}
	frameLimit := sandboxFrameLimit(s.memoryLimitBytes)
	if err := sandboxwire.Write(stdin, frameLimit, request); err != nil {
		stopSandboxWorker(cmd)
		return object.EmeraldValue{}, "", &SandboxWorkerError{Detail: "send worker request", Cause: err}
	}
	if err := stdin.Close(); err != nil {
		stopSandboxWorker(cmd)
		return object.EmeraldValue{}, "", &SandboxWorkerError{Detail: "close worker input", Cause: err}
	}

	phase := ""
	for {
		var response sandboxwire.Response
		if err := sandboxwire.Read(stdout, frameLimit, &response); err != nil {
			stopSandboxWorker(cmd)
			if ctxErr := ctx.Err(); ctxErr != nil {
				return object.EmeraldValue{}, phase, ctxErr
			}
			return object.EmeraldValue{}, phase, sandboxWorkerFailure("read worker response", phase, stderr.String(), err)
		}
		if response.Version != sandboxwire.Version {
			stopSandboxWorker(cmd)
			return object.EmeraldValue{}, phase, sandboxWorkerFailure("worker protocol version mismatch", phase, stderr.String(), nil)
		}
		if response.Phase != "" {
			if response.Phase != sandboxwire.PhaseCompile && response.Phase != sandboxwire.PhaseExecute {
				stopSandboxWorker(cmd)
				return object.EmeraldValue{}, phase, sandboxWorkerFailure("invalid worker phase", phase, stderr.String(), nil)
			}
			if (response.Phase == sandboxwire.PhaseCompile && phase != "") ||
				(response.Phase == sandboxwire.PhaseExecute && phase != string(sandboxwire.PhaseCompile)) {
				stopSandboxWorker(cmd)
				return object.EmeraldValue{}, phase, sandboxWorkerFailure("out-of-order worker phase", phase, stderr.String(), nil)
			}
			if response.Value != nil || response.Error != nil {
				stopSandboxWorker(cmd)
				return object.EmeraldValue{}, phase, sandboxWorkerFailure("invalid worker phase response", phase, stderr.String(), nil)
			}
			phase = string(response.Phase)
			continue
		}
		if (response.Value == nil) == (response.Error == nil) {
			stopSandboxWorker(cmd)
			return object.EmeraldValue{}, phase, sandboxWorkerFailure("invalid worker terminal response", phase, stderr.String(), nil)
		}
		if response.Error != nil && !sandboxErrorPhaseMatches(response.Error.Phase, phase) {
			stopSandboxWorker(cmd)
			return object.EmeraldValue{}, phase, sandboxWorkerFailure("invalid worker error phase", phase, stderr.String(), nil)
		}

		waitErr := cmd.Wait()
		if ctxErr := ctx.Err(); ctxErr != nil {
			return object.EmeraldValue{}, phase, ctxErr
		}
		if response.Error != nil {
			if waitErr != nil {
				return object.EmeraldValue{}, phase, sandboxWorkerFailure("worker exited after returning an error", phase, stderr.String(), waitErr)
			}
			return object.EmeraldValue{}, responseErrorPhase(response.Error, phase), sandboxResponseError(response.Error, stderr.String())
		}
		if waitErr != nil {
			return object.EmeraldValue{}, phase, sandboxWorkerFailure("worker exited after returning a value", phase, stderr.String(), waitErr)
		}

		value, err := decodeSandboxValue(response.Value, frameLimit)
		if err != nil {
			return object.EmeraldValue{}, phase, err
		}
		return value, phase, nil
	}
}

func sandboxErrorPhaseMatches(errorPhase sandboxwire.Phase, phase string) bool {
	if errorPhase == "" {
		return true
	}
	if errorPhase != sandboxwire.PhaseCompile && errorPhase != sandboxwire.PhaseExecute {
		return false
	}
	return string(errorPhase) == phase
}

func stopSandboxWorker(cmd *exec.Cmd) {
	if cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	_ = cmd.Wait()
}

func sandboxWorkerFailure(detail string, phase string, stderr string, cause error) error {
	if stderr != "" {
		detail += ": " + stderr
	}
	return &SandboxWorkerError{Phase: phase, Detail: detail, Cause: cause}
}

func responseErrorPhase(response *sandboxwire.Error, fallback string) string {
	if response.Phase != "" {
		return string(response.Phase)
	}
	return fallback
}

func sandboxResponseError(response *sandboxwire.Error, stderr string) error {
	switch response.Kind {
	case "eval":
		return EvalError{
			ClassName:  response.ClassName,
			Message:    response.Message,
			Stacktrace: append([]byte(nil), response.Stacktrace...),
		}
	case "deadline":
		return context.DeadlineExceeded
	case "unsupported_result":
		return &SandboxUnsupportedResultError{Type: response.ClassName, Reason: response.Message}
	case "worker":
		return sandboxWorkerFailure(response.Message, responseErrorPhase(response, ""), stderr, nil)
	default:
		return sandboxWorkerFailure("invalid worker error response", responseErrorPhase(response, ""), stderr, nil)
	}
}

// ServeSandboxWorker serves one sandbox request on stdin/stdout. It is used by
// the bundled emerald-sandbox-worker executable.
func ServeSandboxWorker(input io.Reader, output io.Writer) (err error) {
	frameLimit := maxSandboxProtocolFrameBytes
	defer func() {
		if recovered := recover(); recovered != nil {
			response := sandboxwire.Response{
				Version: sandboxwire.Version,
				Error: &sandboxwire.Error{
					Kind:       "worker",
					Message:    fmt.Sprintf("worker panic: %v", recovered),
					Stacktrace: runtimedebug.Stack(),
				},
			}
			_ = sandboxwire.Write(output, frameLimit, response)
			err = fmt.Errorf("worker panic: %v", recovered)
		}
	}()

	var request sandboxwire.Request
	if err := sandboxwire.Read(input, frameLimit, &request); err != nil {
		return err
	}
	if request.Version != sandboxwire.Version {
		return writeSandboxWorkerError(output, frameLimit, "worker", "protocol", "unsupported protocol version", "")
	}
	if request.MemoryLimitBytes <= 0 {
		return writeSandboxWorkerError(output, frameLimit, "worker", "request", "memory limit must be positive", "")
	}
	if request.DeadlineUnixNano <= 0 {
		return writeSandboxWorkerError(output, frameLimit, "worker", "request", "deadline is required", "")
	}

	frameLimit = sandboxFrameLimit(request.MemoryLimitBytes)
	// GOMEMLIMIT is set before process startup. The RSS limit is advisory on
	// macOS and deliberately best effort, so failure does not reject an eval.
	_ = applySandboxBestEffortMemoryLimit(request.MemoryLimitBytes)

	ctx, cancel := context.WithDeadline(context.Background(), time.Unix(0, request.DeadlineUnixNano))
	defer cancel()

	var phase sandboxwire.Phase
	var phaseWriteErr error
	onPhase := func(next string) {
		if phaseWriteErr != nil {
			return
		}
		phase = sandboxwire.Phase(next)
		phaseWriteErr = sandboxwire.Write(output, frameLimit, sandboxwire.Response{
			Version: sandboxwire.Version,
			Phase:   phase,
		})
	}

	rt := core.NewRuntime()
	rt.InitSandbox()
	value, evalErr := newEngine(rt).evalFileContextWithPhase(ctx, request.FileName, request.Content, EvalOptions{}, onPhase)
	if phaseWriteErr != nil {
		return phaseWriteErr
	}
	if evalErr != nil {
		return writeSandboxWorkerEvaluationError(output, frameLimit, evalErr, phase)
	}

	encoded, encodeErr := encodeSandboxValue(value)
	if encodeErr != nil {
		var unsupported *SandboxUnsupportedResultError
		if errors.As(encodeErr, &unsupported) {
			return writeSandboxWorkerError(output, frameLimit, "unsupported_result", unsupported.Type, unsupported.Reason, phase)
		}
		return writeSandboxWorkerError(output, frameLimit, "worker", "result", encodeErr.Error(), phase)
	}
	return sandboxwire.Write(output, frameLimit, sandboxwire.Response{
		Version: sandboxwire.Version,
		Value:   encoded,
	})
}

func writeSandboxWorkerEvaluationError(
	output io.Writer,
	frameLimit int,
	err error,
	phase sandboxwire.Phase,
) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return writeSandboxWorkerError(output, frameLimit, "deadline", "", err.Error(), phase)
	}

	var evalErr EvalError
	if errors.As(err, &evalErr) {
		return sandboxwire.Write(output, frameLimit, sandboxwire.Response{
			Version: sandboxwire.Version,
			Error: &sandboxwire.Error{
				Kind:       "eval",
				ClassName:  evalErr.ClassName,
				Message:    evalErr.Message,
				Phase:      phase,
				Stacktrace: append([]byte(nil), evalErr.Stacktrace...),
			},
		})
	}

	return writeSandboxWorkerError(output, frameLimit, "worker", "evaluation", err.Error(), phase)
}

func writeSandboxWorkerError(
	output io.Writer,
	frameLimit int,
	kind string,
	className string,
	message string,
	phase sandboxwire.Phase,
) error {
	return sandboxwire.Write(output, frameLimit, sandboxwire.Response{
		Version: sandboxwire.Version,
		Error: &sandboxwire.Error{
			Kind:      kind,
			ClassName: className,
			Message:   message,
			Phase:     phase,
		},
	})
}

type sandboxLimitedBuffer struct {
	limit int
	mu    sync.Mutex
	buf   strings.Builder
}

func (b *sandboxLimitedBuffer) Write(payload []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	originalLen := len(payload)
	remaining := b.limit - b.buf.Len()
	if remaining > 0 {
		if len(payload) > remaining {
			payload = payload[:remaining]
		}
		_, _ = b.buf.Write(payload)
	}
	return originalLen, nil
}

func (b *sandboxLimitedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	return strings.TrimSpace(b.buf.String())
}
