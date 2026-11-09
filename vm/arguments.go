package vm

import "emerald/object"

// bindPositionalArgs moves the supplied positional arguments into the method's
// parameter slots and reports which parameters received an argument. Ruby gives
// the required parameters the last arguments, so a default fills any gap before
// them, as in `def f(a = 1, b)` called with one argument.
func (vm *VM) bindPositionalArgs(method *object.ClosedBlock, basePointer, given int) []bool {
	numParams := method.NumArgs
	if len(method.DefaultArgSlots) == 0 {
		return nil
	}

	stack := vm.stack()

	supplied := make([]object.EmeraldValue, given)
	copy(supplied, stack[basePointer:basePointer+given])

	hasDefault := make([]bool, numParams)
	for _, slot := range method.DefaultArgSlots {
		hasDefault[slot] = true
	}

	// Optionals take the arguments left after the required parameters are
	// filled from the right; the remaining optionals keep their defaults.
	optionalArguments := given - (numParams - len(method.DefaultArgSlots))

	argGiven := make([]bool, numParams)
	nextArgument := 0

	for slot := 0; slot < numParams; slot++ {
		if hasDefault[slot] && optionalArguments == 0 {
			stack[basePointer+slot] = vm.rt.NULL
			continue
		}

		stack[basePointer+slot] = supplied[nextArgument]
		argGiven[slot] = true
		nextArgument += 1

		if hasDefault[slot] {
			optionalArguments -= 1
		}
	}

	return argGiven
}
