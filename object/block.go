package object

import (
	"emerald/bytecode"
	"fmt"
)

type ExceptionTableEntry struct {
	StartIP            int
	EndIP              int
	HandlerIP          int
	CaughtErrorClasses []string
}

type Block struct {
	*BaseEmeraldValue
	bytecode.Bytecode
	NumLocals      int
	NumArgs        int
	Kwargs         []string
	EnforceArity   bool
	ExceptionTable []ExceptionTableEntry
}

func (b *Block) Class() EmeraldValue          { return nil }
func (b *Block) Super() EmeraldValue          { return nil }
func (b *Block) Ancestors() []EmeraldValue    { return []EmeraldValue{} }
func (b *Block) Type() EmeraldValueType       { return BLOCK_VALUE }
func (b *Block) Inspect() string              { return fmt.Sprintf("#<Block:%p>", b) }
func (b *Block) HashKey() string              { return b.Inspect() }
func (b *Block) SingletonClass() EmeraldValue { return nil }

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
