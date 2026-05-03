# Task: Lazy Singleton Creation

## Description
Currently, every object in Emerald (including every instance of a class) has a singleton class created for it immediately upon instantiation. This is done in `object/class.go`'s `New` method and `object/instance.go`.

```go
func (c *Class) New() *Instance {
	instance := &Instance{}

	singleton := NewSingletonClass(instance, BuiltInMethodSet{}, c)

	instance.class = singleton
	instance.BaseEmeraldValue = singleton.BaseEmeraldValue

	return instance
}
```

This is memory-intensive because most objects never need a singleton class (which is only used for defining per-object methods).

## Goals
- [ ] Modify `object.Instance` and other types to create their singleton class lazily.
- [ ] Ensure that `Class()` method returns the singleton class if it exists, or the base class otherwise.
- [ ] Update `BaseEmeraldValue` or its usage to handle the case where a singleton class hasn't been created yet.
- [ ] Verify that defining a method on a specific instance still works (triggers singleton creation).

## Verification Criteria
- [ ] Memory usage for creating many objects should decrease.
- [ ] Existing tests for singleton classes and instance methods must pass.
