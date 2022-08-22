package types

type Slice[T comparable] struct {
	value []T
}

func NewSlice[T comparable](value ...T) *Slice[T] {
	return &Slice[T]{
		value: value,
	}
}

func (slice Slice[T]) Each(cb func(T)) Slice[T] {
	for _, val := range slice.value {
		cb(val)
	}

	return slice
}

func (slice Slice[T]) Find(cb func(T) bool) *T {
	for _, val := range slice.value {
		if cb(val) {
			return &val
		}
	}

	return nil
}

func (slice Slice[T]) Includes(value T) bool {
	for _, val := range slice.value {
		if val == value {
			return true
		}
	}

	return false
}
