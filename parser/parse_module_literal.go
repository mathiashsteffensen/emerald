package parser

import (
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseModuleLiteral() ast.Expression {
	mod := &ast.ModuleLiteral{Token: p.curToken}

	p.nextToken()

	nameToken := p.curToken

	moduleName, isName := p.parseIdentifierExpression().(ast.IdentifierExpression)
	if !isName {
		p.classOrModuleNameError(nameToken)
		return nil
	}

	mod.Name = moduleName

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	mod.Body = p.parseBlockStatement(lexer.END)

	return mod
}
