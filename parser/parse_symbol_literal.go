package parser

import (
	"emerald/parser/ast"
)

func (p *Parser) parseSymbolLiteral() ast.Expression {
	return &ast.SymbolLiteral{Token: p.curToken, Value: p.curToken.Literal[1:]}
}
