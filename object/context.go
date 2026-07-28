package object

import (
	"context"
	"time"
)

type Context struct {
	Outer                   *Context
	File                    string
	Self                    EmeraldValue
	BreakValue              EmeraldValue
	Block                   EmeraldValue
	ExecutionContext        context.Context
	Yield                   func(kwargs map[string]EmeraldValue, args ...EmeraldValue) EmeraldValue
	BlockGiven              func() bool
	DefaultMethodVisibility MethodVisibility
}

// ExecutionError reports whether execution has been canceled or its deadline has passed.
func (ctx *Context) ExecutionError() error {
	if ctx.ExecutionContext == nil {
		return nil
	}

	if err := ctx.ExecutionContext.Err(); err != nil {
		return err
	}

	if deadline, ok := ctx.ExecutionContext.Deadline(); ok && !time.Now().Before(deadline) {
		return context.DeadlineExceeded
	}

	return nil
}

func (ctx *Context) SetDefaultMethodVisibility(new MethodVisibility) {
	ctx.DefaultMethodVisibility = new
}

func (ctx *Context) ValidateMethodVisibility(receiver EmeraldValue, visibility MethodVisibility, isDefinedOnReceiver bool) bool {
	switch visibility {
	case PRIVATE:
		return ctx.Self == receiver
	case PROTECTED:
		// TODO
		return false
	default:
		return true
	}
}
