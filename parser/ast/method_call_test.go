package ast_test

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"reflect"
	"testing"
)

func TestMethodCallDup(t *testing.T) {
	var dup func(ast.MethodCall) *ast.MethodCall = ast.MethodCall.Dup
	left := &ast.IdentifierExpression{Value: "receiver"}
	argument := &ast.IdentifierExpression{Value: "value"}
	original := ast.MethodCall{
		Left:  left,
		Token: lexer.Token{Literal: "."},
		CallExpression: ast.CallExpression{
			Token:            lexer.Token{Literal: "("},
			Method:           ast.IdentifierExpression{Value: "write="},
			Arguments:        []ast.Expression{argument},
			KeywordArguments: []*ast.HashLiteralElement{{Key: argument, Value: left}},
			Block:            &ast.BlockLiteral{},
			Assignment:       true,
		},
	}
	copy := dup(original)
	if copy == &original || !reflect.DeepEqual(*copy, original) {
		t.Fatal("Dup must return a distinct, equal method call including Assignment")
	}
	copy.Assignment = false
	copy.Method.Value = "other"
	if !original.Assignment || original.Method.Value != "write=" {
		t.Fatal("changing copied scalar fields modified the original")
	}
	if copy.Left != left || copy.Block != original.Block || &copy.Arguments[0] != &original.Arguments[0] || &copy.KeywordArguments[0] != &original.KeywordArguments[0] {
		t.Fatal("Dup must preserve the original shallow-copy behavior")
	}
}
