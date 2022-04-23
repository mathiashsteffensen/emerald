package core

import (
	"emerald/object"
)

func NativeBoolToBooleanObject(input bool) object.EmeraldValue {
	if input {
		return TRUE
	}
	return FALSE
}

func IsError(obj object.EmeraldValue) bool {
	_, ok := obj.(object.EmeraldError)

	return ok
}
