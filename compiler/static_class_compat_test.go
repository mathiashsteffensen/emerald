package compiler_test

import (
	"bytes"
	"emerald/bytecode"
	"emerald/compiler"
	"emerald/core"
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"testing"
)

func TestStaticClassLiteralLegacyReceiver(t *testing.T) {
	legacy := &ast.StaticClassLiteral{Body: &ast.BlockStatement{Statements: []ast.Statement{
		&ast.ExpressionStatement{Expression: &ast.IntegerLiteral{Token: lexer.Token{Literal: "7"}, Value: 7}},
	}}}
	for _, tt := range []struct {
		name      string
		node      *ast.StaticClassLiteral
		op        bytecode.Opcode
		formatted string
	}{
		{"legacy implicit self", legacy, bytecode.OpSelf, "class << self\n\t7\nend"},
		{"explicit self", &ast.StaticClassLiteral{Receiver: &ast.Self{IdentifierExpression: ast.IdentifierExpression{Value: "self"}}, Body: legacy.Body}, bytecode.OpSelf, "class << self\n\t7\nend"},
		{"explicit receiver", &ast.StaticClassLiteral{Receiver: ast.IdentifierExpression{Value: "Object"}, Body: legacy.Body}, bytecode.OpConstantGet, "class << Object\n\t7\nend"},
	} {
		t.Run(tt.name+" compile", func(t *testing.T) {
			rt := core.NewRuntime()
			rt.Init()
			c := compiler.New(nil, rt)
			c.Compile(tt.node)
			first := bytecode.Make(tt.op)
			index := 0
			if tt.op == bytecode.OpConstantGet {
				first = bytecode.Make(tt.op, 0)
				index = 1
			}
			want := bytes.Join([][]byte{first, bytecode.Make(bytecode.OpStaticTrue), bytecode.Make(bytecode.OpPushConstant, index), bytecode.Make(bytecode.OpStaticFalse)}, nil)
			if got := c.Bytecode().Instructions; !bytes.Equal(got, want) {
				t.Fatalf("instructions = %v, want %v (error: %v)", []byte(got), want, c.Err())
			}
		})
		t.Run(tt.name+" format", func(t *testing.T) {
			if got := tt.node.String(0); got != tt.formatted {
				t.Fatalf("format = %q, want %q", got, tt.formatted)
			}
		})
	}
}
