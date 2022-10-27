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

	node.Arguments, node.KeywordArguments = p.parseCallArguments()
	node.Block = p.parseBlockLiteral()

	return node
}

func (p *Parser) parseCallArguments() ([]ast.Expression, []*ast.HashLiteralElement) {
	var (
		args   []ast.Expression
		kwargs []*ast.HashLiteralElement
	)

	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken()
		args, kwargs = p.parseMethodArgsWithParentheses()
	} else {
		args, kwargs = p.parseMethodArgsWithoutParentheses()
	}

	return args, kwargs
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
	lexer.COLON,           // When in a ternary
}

func (p *Parser) peekTokenDoesntSignifyCallArguments() bool {
	return p.peekPrecedence() != LOWEST || p.peekTokenIsMultiple(endOfMethodArgsWithoutParenthesesTokens...)
}

func (p *Parser) parseMethodArgsWithParentheses() ([]ast.Expression, []*ast.HashLiteralElement) {
	if p.peekTokenIs(lexer.RPAREN) {
		p.nextToken()
		return []ast.Expression{}, []*ast.HashLiteralElement{}
	}

	p.nextToken()

	args, kwargs := p.parseArgumentList()

	if !p.expectPeek(lexer.RPAREN) {
		return nil, nil
	}

	return args, kwargs
}

func (p *Parser) parseMethodArgsWithoutParentheses() ([]ast.Expression, []*ast.HashLiteralElement) {
	if p.peekTokenDoesntSignifyCallArguments() {
		return []ast.Expression{}, []*ast.HashLiteralElement{}
	}

	p.nextToken()

	args, kwargs := p.parseArgumentList()

	if p.peekTokenIsMultiple(lexer.NEWLINE, lexer.SEMICOLON) {
		p.nextToken()
	}

	return args, kwargs
}

func (p *Parser) parseArgumentList() ([]ast.Expression, []*ast.HashLiteralElement) {
	startedKeywordArgs := false
	args := []ast.Expression{}
	keywordArgs := []*ast.HashLiteralElement{}

	if p.peekTokenIsMultiple(lexer.COLON, lexer.ARROW) {
		key := p.parseHashLiteralKey()
		p.nextToken()
		keywordArgs = append(keywordArgs, &ast.HashLiteralElement{
			Key:   key,
			Value: p.parseExpression(LOWEST),
		})
		startedKeywordArgs = true
	} else {
		args = append(args, p.parseExpression(LOWEST))
	}

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()

		if startedKeywordArgs {
			key := p.parseHashLiteralKey()
			p.nextToken()
			keywordArgs = append(keywordArgs, &ast.HashLiteralElement{
				Key:   key,
				Value: p.parseExpression(LOWEST),
			})
		} else if p.peekTokenIsMultiple(lexer.COLON, lexer.ARROW) {
			key := p.parseHashLiteralKey()
			p.nextToken()
			keywordArgs = append(keywordArgs, &ast.HashLiteralElement{
				Key:   key,
				Value: p.parseExpression(LOWEST),
			})
			startedKeywordArgs = true
		} else {
			p.nextIfCurNewline()
			args = append(args, p.parseExpression(LOWEST))
		}
	}

	return args, keywordArgs
}
