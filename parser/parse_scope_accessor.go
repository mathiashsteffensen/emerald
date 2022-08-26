package parser

import (
	"emerald/ast"
)

func (p *Parser) parseScopeAccessor(receiver ast.Expression) ast.Expression {
	node := &ast.ScopeAccessor{
		Left:           receiver,
		Token:          p.curToken,
		CallExpression: &ast.CallExpression{},
	}

	p.nextToken()

	node.Method = &ast.IdentifierExpression{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	return node
}
