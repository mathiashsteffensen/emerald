package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseSymbolLiteral() ast.Expression {
	symbol := &ast.SymbolLiteral{Token: p.curToken}

	if !p.expectPeekMultiple(true, lexer.IDENT, lexer.STRING) {
		return nil
	}

	symbol.Value = p.curToken.Literal

	return symbol
}
