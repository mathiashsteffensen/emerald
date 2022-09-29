package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseIndexAccessor(left ast.Expression) ast.Expression {
	node := &ast.MethodCall{Token: p.curToken, Left: left, CallExpression: ast.CallExpression{}}

	node.Arguments = p.parseExpressionList(lexer.RBRACKET)

	node.Method = ast.IdentifierExpression{Value: "[]"}

	return node
}
