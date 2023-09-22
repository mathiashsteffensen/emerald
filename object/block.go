package object

import (
	"emerald/types"
	"fmt"
)

type RescueBlock struct {
	Instructions       []byte
	CaughtErrorClasses *types.Slice[string]
}

func NewRescueBlock(ins []byte, errorClasses ...string) RescueBlock {
	return RescueBlock{
		Instructions:       ins,
		CaughtErrorClasses: types.NewSlice(errorClasses...),
	}
}

type Block struct {
	*BaseEmeraldValue
	Instructions []byte
	NumLocals    int
	NumArgs      int
	Kwargs       []string
	EnforceArity bool
	RescueBlocks []RescueBlock
}

func (b *Block) Class() EmeraldValue       { return nil }
func (b *Block) Super() EmeraldValue       { return nil }
func (b *Block) Ancestors() []EmeraldValue { return []EmeraldValue{} }
func (b *Block) Type() EmeraldValueType    { return BLOCK_VALUE }
func (b *Block) Inspect() string           { return fmt.Sprintf("#<Block:%p>", b) }
func (b *Block) HashKey() string           { return b.Inspect() }

func NewBlock(instructions []byte, numLocals int, numArgs int, kwargs []string, enforceArity bool) *Block {
	return &Block{
		Instructions: instructions,
		NumLocals:    numLocals,
		NumArgs:      numArgs,
		Kwargs:       kwargs,
		EnforceArity: enforceArity,
	}
}

type ClosedBlock struct {
	*Block
	FreeVariables []EmeraldValue
	Context       *Context
	File          string
	Visibility    MethodVisibility
}

func NewClosedBlock(ctx *Context, block *Block, free []EmeraldValue, file string, visibility MethodVisibility) *ClosedBlock {
	return &ClosedBlock{
		Block:         block,
		FreeVariables: free,
		Context:       ctx,
		File:          file,
		Visibility:    visibility,
	}
}
