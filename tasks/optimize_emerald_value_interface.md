# Task: Optimize EmeraldValue Interface

## Description
The `EmeraldValue` interface in `object/emerald_value.go` is currently very large (19+ methods). This leads to several issues:
1. **Memory Overhead**: Every single object (including small ones like Integers, Booleans, Nil) must implement all these methods, usually by embedding `BaseEmeraldValue`.
2. **Allocation**: Small values that could be passed by value or stored more efficiently are often boxed into this large interface.

Ruby implementations like MRI use "Immediate Values" for small integers, booleans, and nil to avoid allocations.

## Goals
- [ ] Redesign the `EmeraldValue` interface to be as small as possible (ideally just `Type()`, `Inspect()`, etc.).
- [ ] Move many of the object-system specific methods (like `NamespaceDefinitionGet`) out of the interface or into a secondary interface.
- [ ] Explore using a struct with a type tag or a smaller interface to reduce the memory footprint of every object.
- [ ] Implement immediate values for:
    - `Integer`
    - `Float` (if possible)
    - `Boolean` (`true`, `false`)
    - `Nil`

## Verification Criteria
- [ ] Significant reduction in memory allocations for simple arithmetic and logical operations.
- [ ] No regression in Ruby-level functionality.
- [ ] Benchmark comparison showing performance improvements.
