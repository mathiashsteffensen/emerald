package core

import (
	"emerald/object"
	"fmt"
)

var StandardError *object.Class

func init() {
	StandardError = object.NewClass("StandardError", Object, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
}

type standardError struct {
	*object.Instance
	Message string
}

func (err *standardError) Inspect() string {
	return fmt.Sprintf("#<%s: %s>", err.ParentClass().(*object.Class).Name, err.Message)
}

func NewStandardError(msg string) object.EmeraldValue {
	return &standardError{StandardError.New(), msg}
}

func IsStandardError(val object.EmeraldValue) bool {
	if super, ok := val.(*object.Class); ok && super.Name == "StandardError" {
		return true
	}

	return false
}
