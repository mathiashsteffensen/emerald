# Task: REPL Persistence and Fix

## Description
The REPL (`iem`) is currently "functionally broken" in terms of persistence and possibly execution. 
Each line entered into the REPL is compiled and executed in a brand-new VM instance with a fresh main frame. This means local variables do not persist between lines.

```ruby
iem(main):001:0> x = 5
=> 5
iem(main):002:0> x
# Error: undefined local variable or method 'x'
```

Additionally, the `machine.Send(evaluated, "inspect", ...)` call in `repl/repl.go` might be causing issues or being redundant if `evaluated.Inspect()` is already available.

## Goals
- [ ] Refactor `repl.Start` to maintain VM state across lines.
- [ ] Ensure local variables defined at the top level persist between lines.
- [ ] Investigate and fix why the REPL is considered "broken" beyond just persistence.
- [ ] Ensure that exceptions in one line don't crash the entire REPL session but are reported correctly.

## Verification Criteria
- [ ] Local variables persist: `x = 5` followed by `x` should return `5`.
- [ ] Multi-line input (already partially supported by the buffer) works correctly with state.
- [ ] REPL doesn't exit on Emerald-level exceptions.
