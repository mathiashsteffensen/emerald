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

		method.Arguments, method.KeywordArguments, method.ArgumentDefaults, method.BlockArgument = p.parseMethodLiteralArguments(lexer.RPAREN, true, true)
	} else if p.peekTokenIs(lexer.IDENT) || p.peekTokenIs(lexer.BIT_AND) {
		method.Arguments, method.KeywordArguments, method.ArgumentDefaults, method.BlockArgument = p.parseMethodLiteralArgumentsWithoutParentheses()
	} else {
		method.Arguments, method.KeywordArguments, method.ArgumentDefaults = []*ast.IdentifierExpression{}, []*ast.IdentifierExpression{}, []ast.Expression{}
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
// a comma the parameter list can continue on the next line, and a trailing
// &block parameter is allowed.
func (p *Parser) parseMethodLiteralArgumentsWithoutParentheses() (args, kwargs []*ast.IdentifierExpression, defaults []ast.Expression, blockArgument *ast.IdentifierExpression) {
	p.nextToken()

	startedKeywordArguments := false

	if p.curTokenIs(lexer.BIT_AND) {
		blockArgument = p.parseBlockArgument()
		if blockArgument == nil {
			return nil, nil, nil, nil
		}
		p.checkDuplicateBlockArgument(blockArgument, args, kwargs)
		if !p.expectMethodParameterListEnd() {
			return nil, nil, nil, nil
		}
		return
	}

	ident := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
	if p.peekTokenIs(lexer.COLON) {
		kwargs = append(kwargs, ident)
		p.nextToken()
		startedKeywordArguments = true
	} else {
		args, defaults = p.addPositionalParameter(args, defaults, ident, true)
	}

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextIfNewline()

		if p.peekTokenIs(lexer.BIT_AND) {
			p.nextToken()
			blockArgument = p.parseBlockArgument()
			if blockArgument == nil {
				return nil, nil, nil, nil
			}
			p.checkDuplicateBlockArgument(blockArgument, args, kwargs)
			if !p.expectMethodParameterListEnd() {
				return nil, nil, nil, nil
			}
			return
		}

		if !p.expectPeek(lexer.IDENT) {
			return nil, nil, nil, nil
		}

		ident := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}

		if startedKeywordArguments {
			if !p.expectPeek(lexer.COLON) {
				return nil, nil, nil, nil
			}

			kwargs = append(kwargs, ident)
		} else if p.peekTokenIs(lexer.COLON) {
			kwargs = append(kwargs, ident)
			p.nextToken()
			startedKeywordArguments = true
		} else {
			args, defaults = p.addPositionalParameter(args, defaults, ident, true)
		}
	}

	if !p.expectMethodParameterListEnd() {
		return nil, nil, nil, nil
	}

	return
}

// expectMethodParameterListEnd checks that the parenthesis-free parameter list
// ends at the next token, and records an error when it does not, as in
// `def method a b`.
func (p *Parser) expectMethodParameterListEnd() bool {
	if p.peekTokenIsMultiple(lexer.NEWLINE, lexer.SEMICOLON, lexer.EOF) {
		return true
	}

	p.methodParametersBoundaryError()
	return false
}

// methodParametersBoundaryError reports the token that cannot continue a
// parameter list without parentheses.
func (p *Parser) methodParametersBoundaryError() {
	msg := fmt.Sprintf(
		"expected a newline, ';' or end of input after the method parameters at line %d, column %d\n%s",
		p.peekToken.Line,
		p.peekToken.Column,
		p.l.Snapshot(p.peekToken),
	)
	p.addError(msg)
}

func (p *Parser) parseMethodLiteralArguments(delim lexer.TokenType, allowDefaults bool, allowBlockArgument bool) (args, kwargs []*ast.IdentifierExpression, defaults []ast.Expression, blockArgument *ast.IdentifierExpression) {
	if p.peekTokenIs(delim) {
		p.nextToken()
		return
	}

	startedKeywordArguments := false

	if allowBlockArgument && p.peekTokenIs(lexer.BIT_AND) {
		blockArgument = p.parseBlockArgumentWithDelimiter(delim)
		p.checkDuplicateBlockArgument(blockArgument, args, kwargs)
		return
	}

	if !p.expectPeek(lexer.IDENT) {
		return nil, nil, nil, nil
	}

	ident := &ast.IdentifierExpression{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if p.peekTokenIs(lexer.COLON) {
		kwargs = append(kwargs, ident)
		p.nextToken()
		startedKeywordArguments = true
	} else {
		args, defaults = p.addPositionalParameter(args, defaults, ident, allowDefaults)
	}

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextIfNewline()

		if p.peekTokenIs(delim) {
			p.nextToken()
			return
		}

		if allowBlockArgument && p.peekTokenIs(lexer.BIT_AND) {
			blockArgument = p.parseBlockArgumentWithDelimiter(delim)
			p.checkDuplicateBlockArgument(blockArgument, args, kwargs)
			return
		}

		if !p.expectPeek(lexer.IDENT) {
			return nil, nil, nil, nil
		}

		ident := &ast.IdentifierExpression{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}

		if startedKeywordArguments {
			if !p.expectPeek(lexer.COLON) {
				return nil, nil, nil, nil
			}

			kwargs = append(kwargs, ident)
		} else if p.peekTokenIsMultiple(lexer.COLON) {
			kwargs = append(kwargs, ident)
			p.nextToken()
			startedKeywordArguments = true
		} else {
			args, defaults = p.addPositionalParameter(args, defaults, ident, allowDefaults)
		}
	}

	if !p.expectPeek(delim) {
		return nil, nil, nil, nil
	}

	return
}

// parseBlockArgument parses the name of a method's &block parameter. The
// current token is '&'.
func (p *Parser) parseBlockArgument() *ast.IdentifierExpression {
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	return &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
}

// checkDuplicateBlockArgument reports a &block parameter whose name repeats a
// positional or keyword parameter.
func (p *Parser) checkDuplicateBlockArgument(blockArgument *ast.IdentifierExpression, args, kwargs []*ast.IdentifierExpression) {
	if blockArgument == nil {
		return
	}

	for _, parameter := range args {
		if parameter.Value == blockArgument.Value {
			p.duplicatedArgumentNameError(blockArgument.Token)
			return
		}
	}

	for _, parameter := range kwargs {
		if parameter.Value == blockArgument.Value {
			p.duplicatedArgumentNameError(blockArgument.Token)
			return
		}
	}
}

// duplicatedArgumentNameError reports a parameter name that the method already
// declares.
func (p *Parser) duplicatedArgumentNameError(token lexer.Token) {
	msg := fmt.Sprintf(
		"duplicated argument name at line %d, column %d\n%s",
		token.Line,
		token.Column,
		p.l.Snapshot(token),
	)
	p.addError(msg)
}

// parseBlockArgumentWithDelimiter parses `&block` and the closing delimiter of
// a parenthesized parameter list, as in `def method(a, &block)`. The next
// token is '&'.
func (p *Parser) parseBlockArgumentWithDelimiter(delim lexer.TokenType) *ast.IdentifierExpression {
	p.nextToken()

	blockArgument := p.parseBlockArgument()
	if blockArgument == nil || !p.expectPeek(delim) {
		return nil
	}

	return blockArgument
}

// addPositionalParameter records a positional parameter and, when allowDefaults
// is set, the default value expression that follows it, if any. The defaults
// stay aligned with the parameters: a required parameter records nil.
func (p *Parser) addPositionalParameter(
	args []*ast.IdentifierExpression,
	defaults []ast.Expression,
	ident *ast.IdentifierExpression,
	allowDefaults bool,
) ([]*ast.IdentifierExpression, []ast.Expression) {
	var defaultValue ast.Expression
	if allowDefaults {
		defaultValue = p.parseParameterDefault()
	}

	return append(args, ident), append(defaults, defaultValue)
}

// parseParameterDefault parses the value expression after a parameter's "=".
// It returns nil when the parameter has no default.
func (p *Parser) parseParameterDefault() ast.Expression {
	if !p.peekTokenIs(lexer.ASSIGN) {
		return nil
	}

	p.nextToken() // the '='
	p.nextToken() // the first token of the default value

	return p.parseExpression(ASSIGN)
}
