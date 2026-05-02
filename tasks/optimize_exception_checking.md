# Task: Optimize Exception Checking in VM Loop

## Context
Currently, the VM checks for exceptions dynamically on every single instruction execution inside `vm.runWhile`:
```go
if !vm.currentFiber().inRescue && vm.ExceptionIsRaised() {
    // ...
}
```
This branch in the hottest loop of the application limits the execution speed of the bytecode interpreter.

## Action Items
- [ ] Remove the continuous `ExceptionIsRaised()` check from the primary `fetch`/`decode`/`execute` cycle.
- [ ] Implement an Exception Table (mapping Instruction Pointer ranges to rescue blocks).
- [ ] Modify the `raise` logic so that it immediately jumps the IP to the correct rescue handler, avoiding the per-instruction overhead.
