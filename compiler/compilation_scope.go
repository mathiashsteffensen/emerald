package compiler

import bytecode2 "emerald/bytecode"

type CompilationScope struct {
	bytecode            bytecode2.Bytecode
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	block               bool   // a non-method block, where `next` returns from the block
	loops               []loop // enclosing while loops, innermost last
}

// loop is a while loop being compiled
type loop struct {
	conditionPosition int   // where `next` jumps to
	breakJumps        []int // `break` jumps to patch to the end of the loop
}
