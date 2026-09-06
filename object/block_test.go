package object

import (
	"emerald/bytecode"
	"strings"
	"testing"
)

func TestBlock(t *testing.T) {
	instructions := bytecode.Instructions{1, 2, 3}
	block := NewBlock(bytecode.Bytecode{Instructions: instructions}, 2, 1, []string{"kw"}, true)

	if block.Type() != BLOCK_VALUE {
		t.Errorf("expected BLOCK_VALUE, got %s", block.Type().String())
	}

	if !strings.HasPrefix(block.Inspect(), "#<Block:0x") {
		t.Errorf("unexpected inspect format: %s", block.Inspect())
	}

	if !block.Class().IsNil() {
		t.Error("block class should be nil")
	}

	if !block.Super().IsNil() {
		t.Error("block super should be nil")
	}

	if len(block.Ancestors()) != 0 {
		t.Error("block should have no ancestors")
	}

	if !block.SingletonClass().IsNil() {
		t.Error("block should have no singleton class")
	}

	if block.HashKey() == "" {
		t.Error("block hash key should not be empty")
	}
}

func TestClosedBlock(t *testing.T) {
	block := &Block{}
	ctx := &Context{}
	closed := NewClosedBlock(ctx, block, nil, "file.rb", PUBLIC)

	if closed.Block != block {
		t.Error("closed block should wrap the original block")
	}

	if closed.Context != ctx {
		t.Error("closed block should have the context")
	}

	if closed.File != "file.rb" {
		t.Error("closed block should have the file path")
	}
}
