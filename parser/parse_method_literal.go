package parser

import (
	"emerald/ast"
	"emerald/lexer"
)

func (p *Parser) parseMethodLiteral() ast.Expression {
	method := &ast.MethodLiteral{Token: p.curToken, BlockLiteral: &ast.BlockLiteral{}}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	method.Name = p.parseIdentifierExpression()
	p.nextIfSemicolonOrNewline()

	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken()

		method.Parameters = p.parseCallArguments()
	} else {
		method.Parameters = make([]ast.Expression, 0)
	}

	p.nextIfSemicolonOrNewline()

	body, rescueBlocks, ensure := p.parseBlockBody()

	method.Body = body
	method.RescueBlocks = rescueBlocks
	method.EnsureBlock = ensure

	return method
}
