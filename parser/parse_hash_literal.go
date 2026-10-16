package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
	"strings"
)

func (p *Parser) parseHashLiteral() ast.Expression {
	token := p.curToken
	values := []*ast.HashLiteralElement{}

	p.nextIfNewline()

	for p.curToken.Type != lexer.RBRACE {
		if p.peekTokenIs(lexer.RBRACE) {
			p.nextToken()
			break
		}

		p.nextToken()

		element := p.parseHashLiteralElement()
		if element == nil {
			return nil
		}

		values = append(values, element)

		if !p.peekTokenIs(lexer.COMMA) {
			p.nextIfNewline()
			if !p.expectPeek(lexer.RBRACE) {
				return nil
			}
		} else {
			p.nextToken()
			p.nextIfNewline()
		}
	}

	return &ast.HashLiteral{
		Values: values,
		Token:  token,
	}
}

// parseHashLiteralElement parses `key: value`, `key => value` or the shorthand `key:`,
// which takes its value from the variable, method or constant named key.
func (p *Parser) parseHashLiteralElement() *ast.HashLiteralElement {
	label := p.curToken
	isLabel := p.peekTokenIs(lexer.COLON)

	key := p.parseHashLiteralKey()
	if key == nil {
		return nil
	}

	if isLabel {
		// Like Ruby, newlines after a label are ignored, so the value may start on the next line
		p.nextIfNewline()
	}

	var value ast.Expression

	// The value is omitted when nothing can start it, e.g. `{key:}`, `f(key:, other: 1)` or `f key:` before `end`
	if isLabel && p.prefixParseFns[p.peekToken.Type] == nil {
		value = p.parseOmittedHashValue(label)
	} else {
		p.nextToken()
		value = p.parseExpression(LOWEST)
	}

	if value == nil {
		return nil
	}

	return &ast.HashLiteralElement{Key: key, Value: value}
}

func (p *Parser) parseHashLiteralKey() ast.Expression {
	var key ast.Expression

	if p.peekTokenIs(lexer.COLON) {
		key = &ast.SymbolLiteral{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken()
	} else {
		key = p.parseExpression(LOWEST)
		if !p.expectPeek(lexer.ARROW) {
			return nil
		}
	}

	return key
}

// parseOmittedHashValue resolves `key:` to `key: key`. Like Ruby, keywords such as `class:` resolve to method calls
func (p *Parser) parseOmittedHashValue(label lexer.Token) ast.Expression {
	if (label.Type != lexer.IDENT && !lexer.IsKeyword(label)) || strings.HasSuffix(label.Literal, "?") || strings.HasSuffix(label.Literal, "!") {
		p.addError(fmt.Sprintf(
			"expected a value for hash key %q at line %d column %d\n%s",
			label.Literal, label.Line, label.Column, p.l.Snapshot(label),
		))
		return nil
	}

	return ast.IdentifierExpression{Token: label, Value: label.Literal}
}
