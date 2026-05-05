# Detailed Plan: Optimize `EmeraldValue` Interface

## Objective
Transition `EmeraldValue` from a large Go interface to a concrete universal `Value` struct. This aims to eliminate interface boxing allocations for immediate values (Integer, Float, Boolean, Nil) and heavily reduce memory overhead, utilizing a NaN-boxing/Word-boxing-like technique.

## Current State Analysis
Currently, `EmeraldValue` is a large interface (19+ methods). All objects, including small primitives, are heap-allocated structs (e.g., `*IntegerInstance`) that implement this interface. When passing values around, the Go runtime allocates heap memory to box these values into interface types. This causes significant performance degradation and memory bloat during primitive operations.

## Proposed Architecture
We will introduce a 24-byte (3-word) concrete struct to represent all values in the runtime.

```go
type EmeraldValueType uint8

const (
    NIL_VALUE EmeraldValueType = iota
    CLASS_VALUE
    INSTANCE_VALUE
    // ... existing types ...
    INTEGER_VALUE
    FLOAT_VALUE
    TRUE_VALUE
    FALSE_VALUE
    NIL_VALUE
)

// EmeraldValue struct replaces the interface
type EmeraldValue struct {
    Type EmeraldValueType
    Ptr  unsafe.Pointer // Points to the heap-allocated object (or the Class for immediate values)
    Num  uint64         // Stores payload for immediate values
}
```

By passing `EmeraldValue` by value everywhere, we eliminate boxing allocations. Immediate values will not allocate any heap memory, and their values will be stored directly in `Num`. For these immediate values, `Ptr` will store a pointer to their corresponding Class `*Class` so that `Class()` and method extractions can resolve instantly without requiring a `*Runtime` context.

## Step-by-Step Implementation Plan

### Phase 1: Preparation & Renaming
1. **Benchmark Allocations**:
   Write and run benchmarks specifically targeting loops with heavy integer arithmetic. Note `allocs/op` values and runtime to a file for later analysis.
2. **Rename the existing interface**:
   Rename the current `EmeraldValue` interface in `object/emerald_value.go` to `HeapObject`. This interface will internally handle method delegations for non-immediate, heap-allocated objects.
3. **Define `EmeraldValue` Struct**:
   Create the new concrete `EmeraldValue` struct. Define `NIL_VALUE` (value `0`) which will act as the new `nil` equivalent for uninitialized or missing values.
4. **Add Immediate Value Types**:
   Add `INTEGER_VALUE`, `FLOAT_VALUE`, `TRUE_VALUE`, `FALSE_VALUE`, `NIL_VALUE` to the `EmeraldValueType` enum.

### Phase 2: Struct Methods & Delegation
1. **Implement `HeapObject` Unboxing**:
   Implement a helper on the new struct, e.g., `func (val EmeraldValue) heap() HeapObject`, which converts the `val.Ptr` into the `HeapObject` interface if the `Type` represents a heap object (like `INSTANCE_VALUE`, `CLASS_VALUE`).
2. **Implement the 19+ Interface Methods**:
   Re-implement all previous interface methods directly on the `EmeraldValue` struct.
   - For immediate values: Handle calls like `Type()`, `Inspect()`, `Class()` internally by inspecting the `Type`, `Num`, and `Ptr`.
   - For mutations like `InstanceVariableSet` or `DefineMethod`: Panic or return an error if called on immediate values.
   - For heap objects: Delegate execution by returning `val.physical().[MethodName](...)`.

### Phase 3: Immediate Values Implementation
1. **Refactor `Integer`**:
   - Remove `IntegerInstance` entirely.
   - Update `rt.NewInteger(val int64)` to return `EmeraldValue{Type: INTEGER_VALUE, Ptr: unsafe.Pointer(rt.Integer), Num: uint64(val)}`.
   - Update integer operators (`integerAdd`, `integerSubtract`, etc.) to cast and extract `val.Num` instead of unboxing an `*IntegerInstance`.
2. **Refactor `Float`**:
   - Apply the same pattern, but use `math.Float64bits()` and `math.Float64frombits()` to losslessly store the 64-bit float in the `Num` field.
3. **Refactor `Boolean` & `Nil`**:
   - Update `rt.TRUE`, `rt.FALSE`, and `rt.NULL` to be globally accessible `EmeraldValue` structs (e.g., `EmeraldValue{Type: TRUE_VALUE}`).

### Phase 4: Widespread Codebase Refactor (Fixing nil)
1. **Fix `nil` checks**:
   Search across `core/`, `vm/`, `compiler/`, and `object/` for code that compares an `EmeraldValue` to `nil` (`if val == nil`). Replace these checks with `if val.Type == object.NIL_VALUE` or a new `val.IsDefined()` method.
2. **Fix Return Types**:
   Update methods that return `nil` for an `EmeraldValue` to return `object.EmeraldValue{}`.

### Phase 5: Verification & Benchmarking
1. **Test Suite Verification**:
   Run the full existing test suite (`go test ./...`) to guarantee Ruby compatibility and runtime semantic equivalence.
2. **Benchmark Allocations**:
   Run benchmarks specifically targeting loops with heavy integer arithmetic. Ensure that `allocs/op` drops significantly, preferably to `0` for pure integer operations. Compare with earlier saved file.

## Success Criteria Checklist
- [ ] `EmeraldValue` is a struct, completely replacing the previous interface.
- [ ] Integer, Float, Boolean, and Nil are implemented as immediate values.
- [ ] Heap allocations (`new(...)`) are avoided entirely when creating immediate values.
- [ ] No compilation errors across the AST, compiler, VM, and Core libraries.
- [ ] Test suites pass completely with no regressions.
- [ ] Benchmarking confirms the memory footprint reduction.
