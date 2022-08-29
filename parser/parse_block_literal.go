package parser

import (
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseBlockLiteral() *ast.BlockLiteral {
	var endToken lexer.TokenType
	if p.peekTokenIs(lexer.LBRACE) {
		endToken = lexer.RBRACE
	} else {
		if p.peekTokenIs(lexer.DO) {
			endToken = lexer.END
		} else {
			return nil
		}
	}

	p.nextToken()

	block := &ast.BlockLiteral{Body: &ast.BlockStatement{}, Token: p.curToken}

	if p.peekTokenIs(lexer.BIT_OR) {
		p.nextToken()
		block.Parameters = p.parseExpressionList(lexer.BIT_OR)
	}

	if endToken == lexer.RBRACE {
		block.Body = p.parseBlockStatement(endToken)
	} else {
		body, rescues, ensure := p.parseBlockBody()
		block.Body = body
		block.RescueBlocks = rescues
		block.EnsureBlock = ensure
	}

	return block
}
