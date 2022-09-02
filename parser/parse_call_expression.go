package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseCallExpression(method ast.Expression) ast.Expression {
	exp := ast.CallExpression{Token: p.curToken, Method: method.(ast.IdentifierExpression)}

	exp.Arguments = p.parseExpressionList(lexer.RPAREN)

	exp.Block = p.parseBlockLiteral()

	return exp
}
