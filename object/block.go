package object

import (
	"emerald/bytecode"
	"emerald/types"
	"fmt"
)

type RescueBlock struct {
	bytecode.Bytecode
	CaughtErrorClasses *types.Slice[string]
}

func NewRescueBlock(bytecode bytecode.Bytecode, errorClasses ...string) RescueBlock {
	return RescueBlock{
		Bytecode:           bytecode,
		CaughtErrorClasses: types.NewSlice(errorClasses...),
	}
}

type Block struct {
	*BaseEmeraldValue
	bytecode.Bytecode
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

func NewBlock(bytecode bytecode.Bytecode, numLocals int, numArgs int, kwargs []string, enforceArity bool) *Block {
	return &Block{
		Bytecode:     bytecode,
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
