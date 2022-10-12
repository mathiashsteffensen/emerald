package parser

import (
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseIdentifierExpression() ast.Expression {
	node := ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(lexer.LBRACE) || p.peekTokenIs(lexer.DO) {
		callExpression := ast.CallExpression{
			Token:     p.curToken,
			Method:    node,
			Arguments: []ast.Expression{},
		}

		callExpression.Block = p.parseBlockLiteral()

		return callExpression
	}

	if p.peekTokenIs(lexer.ASSIGN) {
		p.nextToken()
		return p.parseAssignmentExpression(node)
	}

	if p.peekTokenDoesntSignifyCallArguments() {
		return node
	}

	callExpression := ast.CallExpression{
		Token:  p.curToken,
		Method: node,
	}

	callExpression.Arguments = p.parseMethodArgsWithoutParentheses()
	callExpression.Block = p.parseBlockLiteral()

	return callExpression
}

func (p *Parser) parseInstanceVariable() ast.Expression {
	return &ast.InstanceVariable{IdentifierExpression: ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}}
}

func (p *Parser) parseGlobalVariable() ast.Expression {
	return &ast.GlobalVariable{IdentifierExpression: ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}}
}

func (p *Parser) parseSelf() ast.Expression {
	return &ast.Self{IdentifierExpression: ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}}
}
