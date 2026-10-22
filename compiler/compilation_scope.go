package compiler

import bytecode2 "emerald/bytecode"

type CompilationScope struct {
	bytecode            bytecode2.Bytecode
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	block               bool  // a non-method block, where `next` returns from the block
	loops               []int // condition positions of enclosing while loops, innermost last
}
