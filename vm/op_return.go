package vm

import "emerald/object"

func (vm *VM) executeOpReturn() {
	vm.currentFiber().popFrame()

	vm.push(vm.rt.NULL)
}

func (vm *VM) executeOpReturnValue() {
	returnValue := vm.stack()[vm.currentFiber().sp-1]

	vm.currentFiber().popFrame()

	vm.push(returnValue)
}

// returnSignal is panicked by OpReturnFromBlock to unwind the Go stack, through any native methods, to its home invocation
type returnSignal struct {
	home  *object.Invocation
	value object.EmeraldValue
}

func (vm *VM) executeOpReturnFromBlock() {
	value := vm.pop()
	home := vm.currentFiber().currentFrame().homeInvocation()

	if home == nil || !home.Running {
		vm.rt.Raise(vm.rt.NewLocalJumpError("unexpected return"))
		return
	}

	panic(returnSignal{home: home, value: value})
}

// runAsHome runs the method invocation in frame, the current frame, while condition holds, so a `return` in a block literal it creates
// ends the invocation with the return value
func (vm *VM) runAsHome(frame *Frame, condition func() bool) {
	home := &object.Invocation{Running: true}
	frame.home = home

	ctx := vm.ctx
	fiber := vm.currentFiber()
	framesIndex := fiber.framesIndex

	defer func() {
		home.Running = false

		if r := recover(); r != nil {
			signal, ok := r.(returnSignal)
			if !ok || signal.home != home {
				panic(r)
			}

			vm.ctx = ctx
			fiber.framesIndex = framesIndex
			fiber.popFrame()
			fiber.push(signal.value)
		}
	}()

	vm.runWhile(condition)
}
