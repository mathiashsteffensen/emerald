# Add crash regression tests

Goal: lock in the hardening work with tests that prove bad Emerald code does not crash Go.

Scope:

- Add tests for invalid method args.
- Add tests for missing files and bad `require_relative` input.
- Add tests for bad regexp input.
- Add tests for wrong receiver/type in core methods.
- Add tests for exceptions raised inside blocks.
- Add tests for web-handler-like repeated execution if relevant.

Test style:

- Prefer normal Go tests around compiler/VM/core.
- Use `defer recover` only where verifying "does not panic" directly helps.
- Assert Emerald exceptions or returned errors, not only absence of panic.

Done when:

- Known invalid inputs fail predictably.
- `EM_TEST=1 go test ./...` covers the new cases.
