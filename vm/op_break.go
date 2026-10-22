package vm

import "emerald/object"

// breakTarget is a method call given a block literal. A `break` in that block returns from the call.
type breakTarget struct {
	block *object.ClosedBlock
}

// breakSignal is panicked by OpBreak to unwind the Go stack, through any native methods, to its target call
type breakSignal struct {
	target *breakTarget
	value  object.EmeraldValue
}

func (vm *VM) executeOpBreak() {
	value := vm.pop()
	block := vm.currentFiber().currentFrame().block

	// The outermost match is the call the block literal was given to, inner ones received it passed along
	for _, target := range vm.breakTargets {
		if target.block == block {
			panic(breakSignal{target: target, value: value})
		}
	}

	vm.rt.Raise(vm.rt.NewLocalJumpError("break from proc-closure"))
}

// callWithBreakTarget runs call so a `break` in block returns from it.
// The break value replaces everything on the stack from receiverPos, where the call's receiver is.
func (vm *VM) callWithBreakTarget(block *object.ClosedBlock, receiverPos int, call func()) {
	target := &breakTarget{block: block}
	vm.breakTargets = append(vm.breakTargets, target)

	ctx := vm.ctx
	fiber := vm.currentFiber()
	framesIndex := fiber.framesIndex

	defer func() {
		vm.breakTargets = vm.breakTargets[:len(vm.breakTargets)-1]

		if r := recover(); r != nil {
			signal, ok := r.(breakSignal)
			if !ok || signal.target != target {
				panic(r)
			}

			vm.ctx = ctx
			fiber.framesIndex = framesIndex
			fiber.sp = receiverPos
			fiber.push(signal.value)
		}
	}()

	call()
}
