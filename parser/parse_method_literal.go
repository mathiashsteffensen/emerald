package parser

import (
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
)

func (p *Parser) parseMethodLiteral() ast.Expression {
	method := &ast.MethodLiteral{Token: p.curToken, BlockLiteral: &ast.BlockLiteral{}}

	p.nextToken()

	methodIdent := ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
	if p.curTokenIs(lexer.LBRACKET) {
		if !p.expectPeek(lexer.RBRACKET) {
			return nil
		}
		methodIdent.Value = "[]"
		methodIdent.Token.Literal = "[]"
	}

	if p.peekTokenIs(lexer.ASSIGN) {
		methodIdent.Value = methodIdent.Value + p.peekToken.Literal
		methodIdent.Token.Literal = methodIdent.TokenLiteral() + p.peekToken.Literal
		p.nextToken()
	}

	method.Name = methodIdent

	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken()

		method.Arguments, method.KeywordArguments = p.parseMethodLiteralArguments(lexer.RPAREN)
	} else if p.peekTokenIs(lexer.IDENT) {
		method.Arguments, method.KeywordArguments = p.parseMethodLiteralArgumentsWithoutParentheses()
	} else {
		method.Arguments, method.KeywordArguments = []*ast.IdentifierExpression{}, []*ast.IdentifierExpression{}
	}

	p.nextIfSemicolonOrNewline()

	body, rescueBlocks, ensure := p.parseBlockBody()

	method.Body = body
	method.RescueBlocks = rescueBlocks
	method.EnsureBlock = ensure

	return method
}

// parseMethodLiteralArgumentsWithoutParentheses parses the parameters that
// follow a method name without parentheses, as in `def greet name, age`. After
// a comma the parameter list can continue on the next line.
func (p *Parser) parseMethodLiteralArgumentsWithoutParentheses() (args []*ast.IdentifierExpression, kwargs []*ast.IdentifierExpression) {
	p.nextToken()

	startedKeywordArguments := false

	ident := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
	if p.peekTokenIs(lexer.COLON) {
		kwargs = append(kwargs, ident)
		p.nextToken()
		startedKeywordArguments = true
	} else {
		args = append(args, ident)
	}

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextIfNewline()

		if !p.expectPeek(lexer.IDENT) {
			return nil, nil
		}

		ident := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}

		if startedKeywordArguments {
			if !p.expectPeek(lexer.COLON) {
				return nil, nil
			}

			kwargs = append(kwargs, ident)
		} else if p.peekTokenIs(lexer.COLON) {
			kwargs = append(kwargs, ident)
			p.nextToken()
			startedKeywordArguments = true
		} else {
			args = append(args, ident)
		}
	}

	if !p.peekTokenIsMultiple(lexer.NEWLINE, lexer.SEMICOLON, lexer.EOF) {
		p.methodParametersBoundaryError()
		return nil, nil
	}

	return
}

// methodParametersBoundaryError reports the token that cannot continue a
// parameter list without parentheses, such as the b in `def method a b`.
func (p *Parser) methodParametersBoundaryError() {
	msg := fmt.Sprintf(
		"expected a newline, ';' or end of input after the method parameters at line %d, column %d\n%s",
		p.peekToken.Line,
		p.peekToken.Column,
		p.l.Snapshot(p.peekToken),
	)
	p.addError(msg)
}

func (p *Parser) parseMethodLiteralArguments(delim lexer.TokenType) (args []*ast.IdentifierExpression, kwargs []*ast.IdentifierExpression) {
	if p.peekTokenIs(delim) {
		p.nextToken()
		return
	}

	startedKeywordArguments := false

	if !p.expectPeek(lexer.IDENT) {
		return nil, nil
	}

	if p.peekTokenIs(lexer.COLON) {
		kwargs = append(kwargs, &ast.IdentifierExpression{
			Token: p.curToken,
			Value: p.curToken.Literal,
		})
		p.nextToken()
		startedKeywordArguments = true
	} else {
		args = append(args, &ast.IdentifierExpression{
			Token: p.curToken,
			Value: p.curToken.Literal,
		})
	}

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextIfNewline()

		if p.peekTokenIs(delim) {
			p.nextToken()
			return
		}

		if !p.expectPeek(lexer.IDENT) {
			return nil, nil
		}

		ident := &ast.IdentifierExpression{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}

		if startedKeywordArguments {
			if !p.expectPeek(lexer.COLON) {
				return nil, nil
			}

			kwargs = append(kwargs, ident)
		} else if p.peekTokenIsMultiple(lexer.COLON) {
			kwargs = append(kwargs, ident)
			p.nextToken()
			startedKeywordArguments = true
		} else {
			args = append(args, ident)
		}
	}

	if !p.expectPeek(delim) {
		return nil, nil
	}

	return
}
