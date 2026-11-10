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
		if p.peekTokenStartsDoBlock() {
			endToken = lexer.END
		} else {
			return nil
		}
	}

	p.nextToken()

	defer p.allowDoBlocks()()

	block := &ast.BlockLiteral{Body: &ast.BlockStatement{}, Token: p.curToken}

	if p.peekTokenIs(lexer.BIT_OR) {
		p.nextToken()
		block.Arguments, block.KeywordArguments, _, _ = p.parseMethodLiteralArguments(lexer.BIT_OR, false, false)
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

func (p *Parser) peekTokenStartsDoBlock() bool {
	return p.peekTokenIs(lexer.DO) && !p.doStartsLoopBody
}
