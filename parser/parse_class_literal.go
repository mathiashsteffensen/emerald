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

	nameToken := p.curToken

	className, isName := p.parseIdentifierExpression().(ast.IdentifierExpression)
	if !isName {
		p.classOrModuleNameError(nameToken)
		return nil
	}

	class.Name = className

	if p.peekTokenIs(lexer.LT) {
		p.nextToken()
		p.nextToken()
		class.Parent = p.parseExpression(LOWEST)
	}

	p.nextIfSemicolonOrNewline()

	class.Body = p.parseBlockStatement(lexer.END)

	return class
}

// classOrModuleNameError reports a class or module name that is not a single
// name, such as the assignment in `class X = 1`, at the token it starts with.
func (p *Parser) classOrModuleNameError(nameToken lexer.Token) {
	msg := fmt.Sprintf(
		"class/module name must be CONSTANT at line %d, column %d\n%s",
		nameToken.Line,
		nameToken.Column,
		p.l.Snapshot(nameToken),
	)
	p.addError(msg)
}
