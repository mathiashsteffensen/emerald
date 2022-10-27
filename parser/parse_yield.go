package parser

import (
	"emerald/parser/ast"
)

func (p *Parser) parseYield() ast.Expression {
	yield := ast.Yield{
		Token:     p.curToken,
		Arguments: []ast.Expression{},
	}

	yield.Arguments, yield.KeywordArguments = p.parseCallArguments()

	return yield
}
