package kernel

import "emerald/object"

const GlobalsSize = 65536

var (
	ConstantPool       []object.EmeraldValue
	GlobalVariablePool []object.EmeraldValue
)

func init() {
	Reset()
}

func GetConst(index uint16) object.EmeraldValue {
	return ConstantPool[index]
}

func AddConst(obj object.EmeraldValue) int {
	ConstantPool = append(ConstantPool, obj)
	return len(ConstantPool) - 1
}

func GetGlobalVariable(index uint16) object.EmeraldValue {
	return GlobalVariablePool[index]
}

func SetGlobalVariable(index uint16, obj object.EmeraldValue) {
	GlobalVariablePool[index] = obj
}

func Reset() {
	ConstantPool = []object.EmeraldValue{}
	GlobalVariablePool = make([]object.EmeraldValue, GlobalsSize)
}
