package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
)

func (p *Parser) parseClassLiteral() ast.Expression {
	if p.peekTokenIs(lexer.APPEND) {
		return p.parseStaticClassLiteral()
	}

	class := &ast.ClassLiteral{Token: p.curToken}

	p.nextToken()

	className := p.parseIdentifierExpression()

	var ok bool
	class.Name, ok = className.(ast.IdentifierExpression)
	if !ok {
		p.addError(fmt.Sprintf("Class name was not identifier, got %T", className))
		return nil
	}

	if p.peekTokenIs(lexer.LT) {
		p.nextToken()
		p.nextToken()
		class.Parent = p.parseExpression(LOWEST)
	}

	p.nextIfSemicolonOrNewline()

	class.Body = p.parseBlockStatement(lexer.END)

	return class
}
