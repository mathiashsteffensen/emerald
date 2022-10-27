package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseClassLiteral() ast.Expression {
	if p.peekTokenIs(lexer.APPEND) {
		return p.parseStaticClassLiteral()
	}

	class := &ast.ClassLiteral{Token: p.curToken}

	p.nextToken()

	className := p.parseIdentifierExpression()

	class.Name = className.(ast.IdentifierExpression)

	if p.peekTokenIs(lexer.LT) {
		p.nextToken()
		p.nextToken()
		class.Parent = p.parseExpression(LOWEST)
	}

	p.nextIfSemicolonOrNewline()

	class.Body = p.parseBlockStatement(lexer.END)

	return class
}
