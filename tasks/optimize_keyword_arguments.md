# Task: Optimize keyword arguments

In BuiltInMethod, forcing a map allocation for keyword arguments on every single method call will crush your garbage collector.

The Fix: Pass keyword arguments as a flat slice or rely on the VM stack.

```go
// Flattened kwargs: [key1, val1, key2, val2]
type BuiltInMethod func(ctx *Context, args []EmeraldValue, kwargs []EmeraldValue) EmeraldValue
```
