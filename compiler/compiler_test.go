package compiler

import (
	"emerald/core"
	"emerald/object"
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Compiling panicked %s", r)
		}
	}()

	rt := core.NewRuntime()
	rt.Init()
	Compile("test.rb", "puts(\"Hello\")", rt)
}

func TestCompileDoesNotEmitBytecodeAfterParserError(t *testing.T) {
	rt := core.NewRuntime()
	rt.Init()

	bc := Compile("test.rb", `
result = {}
$symbols.each do |symbol|
	if ema(symbol, 10) > ema(symbol, 30)
		result[symbol] = { side: BUY, position: { type: :percent, value: 100/$symbols.size }
	else
		result[symbol] = { side: SELL, position: { type: :percent, value: 100 }
	end
end
result
`, rt)

	exception := rt.Heap.GetGlobalVariableString("$!")
	if exception.IsNil() {
		t.Fatal("expected parser exception to be raised")
	}

	message := exception.Heap.(object.EmeraldError).Message()
	if !strings.Contains(message, "failed to parse source file test.rb") {
		t.Fatalf("expected parser exception, got %q", message)
	}
	if !strings.Contains(message, "expected next token to be }, got ELSE instead") {
		t.Fatalf("expected missing hash terminator parse error, got %q", message)
	}
	if len(bc.Instructions) != 0 {
		t.Fatalf("expected compiler to emit no bytecode after parser error, got:\n%s", bc.Instructions)
	}
}

func TestCompilerAbortsOnNilASTNode(t *testing.T) {
	rt := core.NewRuntime()
	rt.Init()

	l := lexer.New(lexer.NewInput("test.rb", ""))
	compiler := New(l, rt)

	compiler.Compile(&ast.ExpressionStatement{
		Token:      lexer.Token{Literal: "nil-child"},
		Expression: nil,
	})

	exception := rt.Heap.GetGlobalVariableString("$!")
	if exception.IsNil() {
		t.Fatal("expected compiler exception to be raised")
	}

	message := exception.Heap.(object.EmeraldError).Message()
	if message != "compiler received nil AST node" {
		t.Fatalf("expected nil AST node compiler error, got %q", message)
	}
	if len(compiler.Bytecode().Instructions) != 0 {
		t.Fatalf("expected compiler to emit no bytecode after nil AST node, got:\n%s", compiler.Bytecode().Instructions)
	}
}
