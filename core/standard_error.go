package core

import (
	"emerald/object"
	"fmt"
)

var StandardError *object.Class

func init() {
	StandardError = object.NewClass("StandardError", Object, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
}

type StandardErrorInstance struct {
	*object.Instance
	message string
}

func (err *StandardErrorInstance) Inspect() string {
	return fmt.Sprintf("#<%s: %s>", err.ParentClass().(*object.Class).Name, err.message)
}

func (err *StandardErrorInstance) Message() string {
	return err.message
}

func NewStandardError(msg string) object.EmeraldValue {
	return &StandardErrorInstance{StandardError.New(), msg}
}

func IsStandardError(val object.EmeraldValue) bool {
	if _, ok := val.(*StandardErrorInstance); ok {
		return true
	}

	return false
}
