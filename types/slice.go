package types

type Slice[T comparable] struct {
	Value []T
}

func NewSlice[T comparable](value ...T) *Slice[T] {
	return &Slice[T]{
		Value: value,
	}
}

func (slice *Slice[T]) Each(cb func(T)) *Slice[T] {
	for _, val := range slice.Value {
		cb(val)
	}

	return slice
}

func (slice *Slice[T]) Map(cb func(element T) T) *Slice[T] {
	newSlice := NewSlice[T]()

	for _, val := range slice.Value {
		newSlice.Push(cb(val))
	}

	return newSlice
}

func (slice *Slice[T]) Find(cb func(T) bool) *T {
	for _, val := range slice.Value {
		if cb(val) {
			return &val
		}
	}

	return nil
}

func (slice *Slice[T]) FindIndex(cb func(T) bool) *int {
	for i, val := range slice.Value {
		if cb(val) {
			return &i
		}
	}

	return nil
}

func (slice *Slice[T]) Includes(value T) bool {
	for _, val := range slice.Value {
		if val == value {
			return true
		}
	}

	return false
}

func (slice *Slice[T]) Push(element T) *Slice[T] {
	slice.Value = append(slice.Value, element)

	return slice
}
