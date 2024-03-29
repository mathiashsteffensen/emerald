package parser

import (
	"emerald/parser/ast"
)

func (p *Parser) parseCallExpression(method ast.Expression) ast.Expression {
	exp := ast.CallExpression{Token: p.curToken, Method: method.(ast.IdentifierExpression)}

	exp.Arguments, exp.KeywordArguments = p.parseMethodArgsWithParentheses()

	exp.Block = p.parseBlockLiteral()

	return exp
}
