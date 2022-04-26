package object

type EmeraldError interface {
	EmeraldValue
	Message() string
	ClassName() string
}
