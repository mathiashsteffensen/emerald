# Add defensive built-in argument validation

Goal: built-in methods should validate arity and types before indexing args or casting heap values.

Scope:

- Find built-ins that read `args[0]`, `args[1]`, or cast `arg.Heap.(*Type)` before `EnforceArity`.
- Add arity checks first.
- Replace unchecked casts with `EnforceArgumentType` or targeted checks.
- Return raised Emerald errors consistently.

High-value areas:

- `core/kernel.go`
- `core/string.go`
- `core/io.go`
- `core/regexp.go`
- `core/file.go`
- `core/array.go`
- `core/hash.go`

Done when:

- Calling built-ins with too few, too many, or wrong-typed args cannot panic.
- Regression tests cover representative invalid calls.
