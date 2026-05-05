package heap

import (
	"emerald/object"
)

const GlobalsSize = 65536

type Heap struct {
	SymbolTable  *SymbolTable
	ConstantPool []object.EmeraldValue
	VariablePool []object.EmeraldValue
}

func NewHeap() *Heap {
	return &Heap{
		SymbolTable:  NewSymbolTable(),
		ConstantPool: []object.EmeraldValue{},
		VariablePool: make([]object.EmeraldValue, GlobalsSize),
	}
}

func (h *Heap) GetConstant(index uint16) object.EmeraldValue {
	return h.ConstantPool[index]
}

func (h *Heap) AddConstant(obj object.EmeraldValue) int {
	h.ConstantPool = append(h.ConstantPool, obj)
	return len(h.ConstantPool) - 1
}

func (h *Heap) GetGlobalVariable(index uint16) object.EmeraldValue {
	return h.VariablePool[index]
}

func (h *Heap) GetGlobalVariableString(name string) object.EmeraldValue {
	if symbol, ok := h.SymbolTable.Resolve(name); !ok {
		return object.EmeraldValue{}
	} else {
		return h.GetGlobalVariable(uint16(symbol.Index))
	}
}

func (h *Heap) SetGlobalVariable(index uint16, obj object.EmeraldValue) {
	h.VariablePool[index] = obj
}

func (h *Heap) SetGlobalVariableString(name string, value object.EmeraldValue) {
	symbol, ok := h.SymbolTable.Resolve(name)
	if !ok {
		symbol = h.SymbolTable.DefineGlobal(name)
	}

	h.SetGlobalVariable(uint16(symbol.Index), value)
}
