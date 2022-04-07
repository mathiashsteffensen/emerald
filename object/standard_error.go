package object

import "fmt"

var StandardError *Class

func init() {
	StandardError = NewClass("StandardError", Object, BuiltInMethodSet{}, BuiltInMethodSet{})
}

type standardError struct {
	*Instance
	Message string
}

func (err *standardError) Inspect() string {
	return fmt.Sprintf("#<%s: %s>", err.class.Name, err.Message)
}

func NewStandardError(msg string) EmeraldValue {
	return &standardError{StandardError.New(), msg}
}

func IsStandardError(val EmeraldValue) bool {
	if super, ok := val.(*Class); ok && super.Name == "StandardError" {
		return true
	}

	return false
}
