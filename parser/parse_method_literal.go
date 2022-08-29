package parser

import (
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseMethodLiteral() ast.Expression {
	method := &ast.MethodLiteral{Token: p.curToken, BlockLiteral: &ast.BlockLiteral{}}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	methodIdent := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
	p.nextIfSemicolonOrNewline()

	if p.peekTokenIs(lexer.ASSIGN) {
		methodIdent.Value = methodIdent.Value + p.peekToken.Literal
		methodIdent.Token.Literal = methodIdent.TokenLiteral() + p.peekToken.Literal
		p.nextToken()
	}

	method.Name = methodIdent

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
