package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseMethodCall(left ast.Expression) ast.Expression {
	node := &ast.MethodCall{Token: p.curToken, Left: left, CallExpression: &ast.CallExpression{}}

	methodIdent := &ast.IdentifierExpression{Value: p.peekToken.Literal, Token: p.peekToken}

	p.nextToken()

	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken()
		node.Arguments = p.parseCallArguments()
	} else if p.peekTokenIs(lexer.ASSIGN) {
		methodIdent.Value = methodIdent.Value + p.peekToken.Literal
		methodIdent.Token.Literal = methodIdent.TokenLiteral() + p.peekToken.Literal
		p.nextToken()
		p.nextToken()
		node.Arguments = []ast.Expression{p.parseExpression(LOWEST)}
	} else {
		node.Arguments = p.parseMethodArgsWithoutParentheses()
	}

	node.Method = methodIdent
	node.Block = p.parseBlockLiteral()

	if p.curTokenIs(lexer.DOT) {
		return p.parseMethodCall(node)
	}

	return node
}
