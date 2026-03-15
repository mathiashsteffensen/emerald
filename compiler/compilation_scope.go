package compiler

import bytecode2 "emerald/bytecode"

type CompilationScope struct {
	bytecode            bytecode2.Bytecode
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}
