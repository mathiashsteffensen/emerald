package core

import "emerald/object"

func NativeBoolToBooleanObject(input bool) object.EmeraldValue {
	if input {
		return TRUE
	}
	return FALSE
}

func IsError(obj object.EmeraldValue) bool {
	for _, value := range obj.Ancestors() {
		class := value.ParentClass()

		if class != nil {
			if class, ok := class.(*object.Class); ok {
				if class.Name == StandardError.Name {
					return true
				}
			}
		}
	}

	return false
}
