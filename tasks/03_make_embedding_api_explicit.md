# Make the embedding API explicit

Goal: make Emerald safer and clearer to embed inside Quantrading.

Scope:

- Keep `Eval` and `EvalFile` as the basic API.
- Add options for runtime inputs that currently come from process globals.
- Avoid implicit host `os.Args` unless explicitly requested.
- Return structured errors where practical, not only formatted `$!` messages.
- Document which runtime state is per-engine and which state is shared.

Possible API shape:

- `Eval(content string)`
- `EvalFile(fileName string, content string)`
- `EvalWithOptions(content string, opts EvalOptions)`
- `EvalFileWithOptions(fileName string, content string, opts EvalOptions)`

Done when:

- A host application can run Emerald without inheriting unrelated process args/state.
- Integration code has a clear error contract.
