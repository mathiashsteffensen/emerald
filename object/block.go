package object

import (
	"fmt"
)

type Block struct {
	*BaseEmeraldValue
	Instructions []byte
	NumLocals    int
}

func (b *Block) ParentClass() EmeraldValue { return nil }
func (b *Block) Ancestors() []EmeraldValue { return []EmeraldValue{} }
func (*Block) Type() EmeraldValueType      { return BLOCK_VALUE }
func (b *Block) Inspect() string           { return fmt.Sprintf("#<Block:%p>", b) }

func NewBlock(instructions []byte, numLocals int) *Block {
	return &Block{
		Instructions: instructions,
		NumLocals:    numLocals,
	}
}

type ClosedBlock struct {
	*Block
	FreeVariables []EmeraldValue
}

func NewClosedBlock(block *Block, free []EmeraldValue) *ClosedBlock {
	return &ClosedBlock{
		Block:         block,
		FreeVariables: free,
	}
}
