package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseHashLiteral() ast.Expression {
	values := []*ast.HashLiteralElement{}

	p.nextIfNewline()

	for p.curToken.Type != lexer.RBRACE {
		if p.peekTokenIs(lexer.RBRACE) {
			p.nextToken()
			break
		}

		p.nextToken()

		key := p.parseHashLiteralKey()

		p.nextToken()

		values = append(values, &ast.HashLiteralElement{
			Key:   key,
			Value: p.parseExpression(LOWEST),
		})

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
		Token:  p.curToken,
	}
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
