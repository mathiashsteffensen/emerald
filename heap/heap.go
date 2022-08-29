// Package heap is not a heap in the traditional sense since we are blessed by the Go GC.
// The package simply contains all objects in use by the Emerald VM that is not currently on the stack, and is not stored on another object
// i.e. All literals and global variables. Although literals and global variables will be pushed onto the stack
// when for example passed as arguments, or assigned to a local variable
package heap

import (
	"emerald/object"
)

const GlobalsSize = 65536

var (
	GlobalSymbolTable  *SymbolTable
	ConstantPool       []object.EmeraldValue
	GlobalVariablePool []object.EmeraldValue
)

func init() {
	Reset()
}

func GetConstant(index uint16) object.EmeraldValue {
	return ConstantPool[index]
}

func AddConstant(obj object.EmeraldValue) int {
	ConstantPool = append(ConstantPool, obj)
	return len(ConstantPool) - 1
}

func GetGlobalVariable(index uint16) object.EmeraldValue {
	return GlobalVariablePool[index]
}

func GetGlobalVariableString(name string) object.EmeraldValue {
	if symbol, ok := GlobalSymbolTable.Resolve(name); !ok {
		return nil
	} else {
		return GetGlobalVariable(uint16(symbol.Index))
	}
}

func SetGlobalVariable(index uint16, obj object.EmeraldValue) {
	GlobalVariablePool[index] = obj
}

func SetGlobalVariableString(name string, value object.EmeraldValue) {
	symbol, ok := GlobalSymbolTable.Resolve(name)
	if !ok {
		symbol = GlobalSymbolTable.DefineGlobal(name)
	}

	SetGlobalVariable(uint16(symbol.Index), value)
}

func Reset() {
	GlobalSymbolTable = NewSymbolTable()
	ConstantPool = []object.EmeraldValue{}
	GlobalVariablePool = make([]object.EmeraldValue, GlobalsSize)
}
