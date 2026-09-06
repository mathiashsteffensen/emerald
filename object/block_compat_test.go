package object_test

import (
	"emerald/object"
	"testing"
)

func TestClosedBlockLegacyFreeVariables(t *testing.T) {
	var constructor func(*object.Context, *object.Block, []object.EmeraldValue, string, object.MethodVisibility) *object.ClosedBlock = object.NewClosedBlock
	value := object.EmeraldValue{TypeID: object.INTEGER_VALUE, Num: 7}
	free := []object.EmeraldValue{value}
	ctx, block := &object.Context{}, &object.Block{}
	closed := constructor(ctx, block, free, "legacy.em", object.PRIVATE)
	if closed.Block != block || closed.Context != ctx || closed.File != "legacy.em" || closed.Visibility != object.PRIVATE {
		t.Fatal("constructor lost block metadata")
	}
	for _, closed := range []*object.ClosedBlock{closed, {Block: block, FreeVariables: free}} {
		var values []object.EmeraldValue = closed.FreeVariables
		if values[0] != value {
			t.Fatal("free variable value was not preserved")
		}
		if closed.FreeVariable(0) != &free[0] {
			t.Fatal("legacy free variable must use the supplied slice as its cell")
		}
		closed.FreeVariables[0].Num = 9
		if closed.FreeVariable(0).Num != 9 {
			t.Fatal("field mutation was not visible through the cell")
		}
		*closed.FreeVariable(0) = value
		if values[0] != value {
			t.Fatal("cell mutation was not visible through the field")
		}
	}
}

func TestClosedBlockSharedFreeVariableCells(t *testing.T) {
	value := object.EmeraldValue{TypeID: object.INTEGER_VALUE, Num: 7}
	cells := []*object.EmeraldValue{&value}
	first := object.NewClosedBlockWithCells(nil, &object.Block{}, cells, "", object.PUBLIC)
	second := object.NewClosedBlockWithCells(nil, &object.Block{}, cells, "", object.PUBLIC)
	first.FreeVariable(0).Num = 9
	if second.FreeVariable(0).Num != 9 || value.Num != 9 {
		t.Fatal("closures did not share the original cell")
	}
	if first.FreeVariables[0].Num != 7 || second.FreeVariables[0].Num != 7 {
		t.Fatal("cell-backed closures should expose capture-time snapshots")
	}
}
