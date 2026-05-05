package object

type EmeraldError interface {
	HeapObject
	Message() string
	ClassName() string
}
