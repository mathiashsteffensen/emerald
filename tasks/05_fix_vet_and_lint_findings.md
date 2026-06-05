# Fix vet and lint findings

Goal: keep the codebase clean enough that tool findings are meaningful.

Current known issue:

- `go vet ./...` reports unreachable code at `core/io.go:62`.

Scope:

- Fix current `go vet` findings.
- Make `go vet ./...` pass.
- Check why `staticcheck ./...` currently reports `warning: "./..." matched no packages`.
- Once staticcheck runs properly, fix real findings or document intentional exceptions.

Done when:

- `go vet ./...` passes.
- Staticcheck either runs successfully or has a documented local setup issue.
- No tool-reported issue hides a real host-crash path.
