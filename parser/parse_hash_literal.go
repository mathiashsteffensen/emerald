package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseHashLiteral() ast.Expression {
	value := make(map[ast.Expression]ast.Expression)

	p.nextIfNewline()

	for p.curToken.Type != lexer.RBRACE {
		if p.peekTokenIs(lexer.RBRACE) {
			p.nextToken()
			break
		}

		p.nextToken()

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

		p.nextToken()

		value[key] = p.parseExpression(LOWEST)

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
		Value: value,
		Token: p.curToken,
	}
}
