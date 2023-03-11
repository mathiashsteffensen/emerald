package object

type Context struct {
	Outer                   *Context
	File                    string
	Self                    EmeraldValue
	BreakValue              EmeraldValue
	Block                   EmeraldValue
	Yield                   func(args ...EmeraldValue) EmeraldValue
	BlockGiven              func() bool
	DefaultMethodVisibility MethodVisibility
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
