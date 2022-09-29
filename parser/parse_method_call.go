package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseMethodCall(left ast.Expression) ast.Expression {
	node := &ast.MethodCall{
		Token: p.curToken,
		Left:  left,
		CallExpression: ast.CallExpression{
			Method: ast.IdentifierExpression{Value: p.peekToken.Literal, Token: p.peekToken},
		},
	}

	p.nextToken()

	if p.peekTokenIs(lexer.ASSIGN) {
		return node
	}

	node.Arguments = p.parseCallArguments()
	node.Block = p.parseBlockLiteral()

	if p.curTokenIs(lexer.DOT) {
		return p.parseMethodCall(node)
	}

	return node
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression

	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken()
		args = p.parseExpressionList(lexer.RPAREN)
	} else {
		args = p.parseMethodArgsWithoutParentheses()
	}

	return args
}

// Tokens that signify the end of a method call without parentheses
var endOfMethodArgsWithoutParenthesesTokens = []lexer.TokenType{
	lexer.EOF,
	lexer.NEWLINE,
	lexer.SEMICOLON,
	lexer.DO,              // When passed a block
	lexer.LBRACE,          // When passed a block
	lexer.RBRACE,          // When last expression in a single line block
	lexer.RBRACKET,        // When last item in array
	lexer.RPAREN,          // When part of a grouped expression
	lexer.COMMA,           // When part of a list i.e. [condition, 2] does not have arguments
	lexer.ASSIGN,          // When an assignment
	lexer.BOOL_OR_ASSIGN,  // When an assignment
	lexer.BOOL_AND_ASSIGN, // When an assignment
	lexer.ARROW,           // When a hash key
	lexer.BIT_OR,          // When last block argument
}

func (p *Parser) peekTokenDoesntSignifyCallArguments() bool {
	return p.peekPrecedence() != LOWEST || p.peekTokenIsMultiple(endOfMethodArgsWithoutParenthesesTokens...)
}

func (p *Parser) parseMethodArgsWithoutParentheses() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenDoesntSignifyCallArguments() {
		return args
	}

	p.nextToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if p.peekTokenIsMultiple(lexer.NEWLINE, lexer.SEMICOLON) {
		p.nextToken()
	}

	return args
}
