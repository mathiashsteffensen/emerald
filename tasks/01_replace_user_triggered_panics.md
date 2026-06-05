# Replace user-triggered panics with Emerald exceptions

Goal: Emerald code should not be able to crash the host process through normal language/runtime behavior.

Scope:

- Audit `panic`, `debug.Fatal`, and `debug.FatalBug` call sites.
- Keep panics only for impossible VM/compiler invariants.
- Convert user-triggerable failures to Emerald exceptions.
- Prioritize file IO, `require_relative`, regexp handling, method dispatch, keyword args, and bad constants.

Useful checks:

- `rg -n "panic\\(|Fatal\\(|FatalBug" . -g '*.go'`
- `go vet ./...`
- `EM_TEST=1 go test ./...`

Done when:

- Invalid Emerald input returns an error or sets an Emerald exception instead of panicking.
- Remaining panics are documented as internal invariants.
