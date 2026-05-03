package compiler

import (
	"emerald/bytecode"
	"emerald/core"
	"emerald/parser/lexer"
	"testing"
)

func TestCompilerScopes(t *testing.T) {
	rt := core.NewRuntime()
	rt.Init()
	compiler := New(lexer.New(lexer.NewInput("scope_test.rb", "")), rt)

	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)
	}

	globalSymbolTable := compiler.symbolTable
	compiler.emit(bytecode.OpMul, lexer.Token{})
	compiler.enterScope()

	if compiler.scopeIndex != 1 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 1)
	}

	compiler.emit(bytecode.OpSub, lexer.Token{})

	if len(compiler.scopes[compiler.scopeIndex].bytecode.Instructions) != 1 {
		t.Errorf("instructions length wrong. got=%d",
			len(compiler.scopes[compiler.scopeIndex].bytecode.Instructions))
	}

	last := compiler.scopes[compiler.scopeIndex].lastInstruction

	if last.Opcode != bytecode.OpSub {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
			last.Opcode, bytecode.OpSub)
	}

	if compiler.symbolTable.Outer != globalSymbolTable {
		t.Errorf("compiler did not enclose symbolTable")
	}

	compiler.leaveScope()

	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d",
			compiler.scopeIndex, 0)
	}

	if compiler.symbolTable != globalSymbolTable {
		t.Errorf("compiler did not restore global symbol table")
	}

	if compiler.symbolTable.Outer != nil {
		t.Errorf("compiler modified global symbol table incorrectly")
	}

	compiler.emit(bytecode.OpAdd, lexer.Token{})
	if len(compiler.scopes[compiler.scopeIndex].bytecode.Instructions) != 2 {
		t.Errorf("instructions length wrong. got=%d",
			len(compiler.scopes[compiler.scopeIndex].bytecode.Instructions))
	}

	last = compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != bytecode.OpAdd {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
			last.Opcode, bytecode.OpAdd)
	}

	previous := compiler.scopes[compiler.scopeIndex].previousInstruction
	if previous.Opcode != bytecode.OpMul {
		t.Errorf("previousInstruction.Opcode wrong. got=%d, want=%d",
			previous.Opcode, bytecode.OpMul)
	}
}
