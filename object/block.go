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

type FreeBinding struct {
	Index int
	Local bool
}

type Block struct {
	*BaseEmeraldValue
	bytecode.Bytecode
	NumLocals      int
	NumArgs        int
	Kwargs         []string
	EnforceArity   bool
	ExceptionTable []ExceptionTableEntry
	FreeBindings   []FreeBinding
}

func (b *Block) Class() EmeraldValue          { return EmeraldValue{} }
func (b *Block) Super() EmeraldValue          { return EmeraldValue{} }
func (b *Block) Ancestors() []EmeraldValue    { return []EmeraldValue{} }
func (b *Block) Type() EmeraldValueType       { return BLOCK_VALUE }
func (b *Block) Inspect() string              { return fmt.Sprintf("#<Block:%p>", b) }
func (b *Block) HashKey() string              { return b.Inspect() }
func (b *Block) SingletonClass() EmeraldValue { return EmeraldValue{} }

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
	// FreeVariables holds legacy mutable captures. For cell-backed closures it is
	// a capture-time snapshot; use FreeVariable to access the shared live value.
	FreeVariables []EmeraldValue
	freeCells     []*EmeraldValue
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

// NewClosedBlockWithCells captures shared mutable cells. FreeVariables exposes
// their initial values only; execution must access the cells through FreeVariable.
func NewClosedBlockWithCells(ctx *Context, block *Block, cells []*EmeraldValue, file string, visibility MethodVisibility) *ClosedBlock {
	free := make([]EmeraldValue, len(cells))
	for i, cell := range cells {
		free[i] = *cell
	}
	closed := NewClosedBlock(ctx, block, free, file, visibility)
	closed.freeCells = cells
	return closed
}

// FreeVariable returns a shared cell, or the legacy slice element itself.
func (b *ClosedBlock) FreeVariable(index int) *EmeraldValue {
	if b.freeCells != nil {
		return b.freeCells[index]
	}
	return &b.FreeVariables[index]
}
