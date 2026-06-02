package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
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

		key := p.parseHashLiteralKey()
		if key == nil {
			return nil
		}

		p.nextToken()

		value := p.parseExpression(LOWEST)
		if value == nil {
			return nil
		}

		values = append(values, &ast.HashLiteralElement{
			Key:   key,
			Value: value,
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
		Token:  token,
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
