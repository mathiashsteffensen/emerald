package parser

import (
	"emerald/ast"
	"emerald/lexer"
)

func (p *Parser) parseModuleLiteral() ast.Expression {
	mod := &ast.ModuleLiteral{Token: p.curToken}

	p.nextToken()

	mod.Name = p.parseIdentifierExpression().(*ast.IdentifierExpression)
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	mod.Body = p.parseBlockStatement(lexer.END)

	return mod
}
