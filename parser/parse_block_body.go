package parser

import (
	"emerald/ast"
	"emerald/lexer"
)

func (p *Parser) parseBlockBody() (*ast.BlockStatement, []*ast.RescueBlock, *ast.EnsureBlock) {
	block := &ast.BlockStatement{Token: p.curToken}

	p.nextToken()

	block.Statements = p.parseBlockMainBodyPart()

	rescues := []*ast.RescueBlock{}

	if p.curTokenIs(lexer.RESCUE) {
		p.nextToken()
		rescues = p.parseBlockRescueParts()
	}

	var ensure *ast.EnsureBlock

	if p.curTokenIs(lexer.ENSURE) {
		ensure = &ast.EnsureBlock{Body: p.parseBlockStatement(lexer.END)}
	}

	return block, rescues, ensure
}

func (p *Parser) parseBlockMainBodyPart() []ast.Statement {
	stmts := []ast.Statement{}

	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.RESCUE) && !p.curTokenIs(lexer.ENSURE) {
		stmt := p.parseStatement()

		if stmt != nil {
			stmts = append(stmts, stmt)
		}

		p.nextToken()
	}

	return stmts
}

func (p *Parser) parseBlockRescueParts() []*ast.RescueBlock {
	rescues := []*ast.RescueBlock{{Body: &ast.BlockStatement{}}}

	rescues[0].Body.Statements = p.parseBlockMainBodyPart()

	if p.curTokenIs(lexer.RESCUE) {
		p.nextToken()

		rescues = append(rescues, p.parseBlockRescueParts()...)
	}

	return rescues
}
