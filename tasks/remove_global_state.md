# Task: Remove Global State to Enable Concurrency and Embedding

## Context
Currently, the `core` and `heap` packages rely heavily on global state:
- `core.Object`, `core.MainObject`
- `heap.ConstantPool`, `heap.GlobalVariablePool`, etc.

This restricts the runtime to a single instance per process, meaning Emerald cannot be safely embedded into a larger Go application handling concurrent scripts.

## Action Items
- [ ] Create a new struct, `Runtime`, to encapsulate all state related to the Emerald runtime.
- [ ] Move `core` built-ins and `heap` pools to be fields on this struct rather than package globals.
- [ ] Update `vm.New()` and related components to initialize and reference this instance-specific state.
- [ ] Add a `emerald.Engine` struct for tieing together the lexer, parser, compiler, runtime and VM. And make sure the API is nice and usable for embedding.
