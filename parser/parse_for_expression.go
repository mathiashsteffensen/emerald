package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
)

func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	expression.Variable = ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(lexer.IN) {
		return nil
	}

	p.nextToken()

	enclosingDoStartsLoopBody := p.doStartsLoopBody
	defer func() { p.doStartsLoopBody = enclosingDoStartsLoopBody }()

	p.doStartsLoopBody = true
	expression.Iterable = p.parseExpression(LOWEST)

	if !p.expectLoopBodyStart() {
		return nil
	}

	p.doStartsLoopBody = false
	expression.Body = p.parseBlockStatement(lexer.END)

	return expression
}

// expectLoopBodyStart moves past the `do`, newline or `;` that ends the loop header
func (p *Parser) expectLoopBodyStart() bool {
	// A command call collection, `for x in m 1, 2`, already moved past its newline or `;`
	if p.curTokenIs(lexer.NEWLINE, lexer.SEMICOLON) {
		return true
	}

	if p.peekTokenIs(lexer.DO) {
		p.nextToken()
		return true
	}

	if p.peekTokenIsMultiple(lexer.NEWLINE, lexer.SEMICOLON) {
		p.nextIfSemicolonOrNewline()
		return true
	}

	if p.peekTokenIs(lexer.EOF) {
		p.unexpectedEofError()
		return false
	}

	p.addError(fmt.Sprintf(
		"expected 'do', a newline or ';' after the for loop collection at line %d, column %d\n%s",
		p.peekToken.Line,
		p.peekToken.Column,
		p.l.Snapshot(p.peekToken),
	))

	return false
}
