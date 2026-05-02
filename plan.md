# Implementation Plan: Remove Global State & Add `emerald.Engine`

## Step 1: Define `core.Runtime` and `heap.Heap`
[X] - **Heap**: Remove global variables (`ConstantPool`, `GlobalVariablePool`, `GlobalSymbolTable`) in `heap` and encapsulate them in a new `heap.Heap` struct.
[X] - **Runtime**: Create a `core.Runtime` (or `emerald.Runtime`) struct. This will serve as the central registry for all language state.
  - It will hold a reference to `*heap.Heap`.
  - It will hold references to all core classes (`ObjectClass`, `StringClass`, `ArrayClass`, etc.) and singletons (`TRUE`, `FALSE`, `NULL`, `MainObject`).
  - It will provide the factory methods currently available globally (e.g., `rt.NewString(...)`, `rt.NewArray(...)`).

## Step 2: Refactor `core` Initialization
[X] - Remove `core/init.go` reliance on Go's `init()`.
[X] - Convert `InitObject()`, `InitString()`, etc., into methods that configure a `core.Runtime` instance during instantiation (e.g., `core.NewRuntime()`).
[X] – Refactor all `core` files to define methods on `*core.Runtime` rather than acting as package-level functions.

## Step 3: Thread the Runtime through Execution (`object.Context`)
- Update `object.Context` to include a reference to `*core.Runtime`. How do we resolve dependencies? Won't they be circular?
- Because built-in methods (which use the `BuiltInMethod` signature) no longer have access to global classes, they will retrieve the `Runtime` from the passed `*object.Context` (e.g., `ctx.Runtime.NewString(...)`) to instantiate new objects or access singletons.

## Step 4: Refactor the Compiler and VM
- The `compiler.Compiler` will need access to `*core.Runtime` (or at least the `*heap.Heap` instance) to register constants and globals locally rather than globally.
- `vm.New()` will be updated to take the compiled bytecode *and* the `*core.Runtime` (or it will create one), ensuring the VM only mutates its isolated instance state.

## Step 5: Implement `emerald.Engine` (Embedding API)
- Create `engine.go` in the root `emerald` package to define the `emerald.Engine` struct.
- `Engine` will tie everything together, holding the `Runtime` and providing a clean execution API for embedding:
  ```go
  engine := emerald.New()
  result, err := engine.Eval("2 + 2")
  ```
- It will wrap the boilerplate of setting up the Lexer -> Parser -> Compiler -> VM pipeline.

## Step 6: Fix Commands and Tests
- Refactor `cmd/emerald/main.go` and `cmd/iem/main.go` to use `emerald.Engine`.
- Refactor all tests across `parser`, `compiler`, `vm`, and `core` to instantiate a `Runtime` or `Engine` instead of relying on the removed global state.
