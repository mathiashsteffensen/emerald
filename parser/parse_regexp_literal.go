package parser

import "emerald/parser/ast"

func (p *Parser) parseRegexpLiteral() ast.Expression {
	return &ast.RegexpLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}
