package parser

import (
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
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

	var ensure *ast.EnsureBlock = nil

	if p.curTokenIs(lexer.ENSURE) {
		ensure = &ast.EnsureBlock{Token: p.curToken, Body: p.parseBlockStatement(lexer.END)}
	}

	return block, rescues, ensure
}

func (p *Parser) parseBlockMainBodyPart() []ast.Statement {
	stmts := []ast.Statement{}

	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.RESCUE) && !p.curTokenIs(lexer.ENSURE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			stmts = append(stmts, stmt)
		}

		p.nextToken()
	}

	if p.curTokenIs(lexer.EOF) {
		p.unexpectedEofError()
		return nil
	}

	return stmts
}

func (p *Parser) parseBlockRescueParts() []*ast.RescueBlock {
	rescues := []*ast.RescueBlock{{Body: &ast.BlockStatement{}}}

	rescues[0].CaughtErrorClasses = p.parseRescueBlockErrorClasses()

	if p.peekTokenIs(lexer.ARROW) {
		p.nextToken()
		p.nextToken()
		rescues[0].ErrorVarName = p.parseIdentifierExpression()
		p.nextToken()
	}

	rescues[0].Body.Statements = p.parseBlockMainBodyPart()

	if p.curTokenIs(lexer.RESCUE) {
		p.nextToken()

		rescues = append(rescues, p.parseBlockRescueParts()...)
	}

	return rescues
}

func (p *Parser) parseRescueBlockErrorClasses() []ast.Expression {
	errorClasses := []ast.Expression{}

	if p.curTokenIs(lexer.ARROW) || p.curTokenIs(lexer.NEWLINE) {
		p.nextToken()
		return errorClasses
	}

	errorClasses = append(errorClasses, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		errorClasses = append(errorClasses, p.parseExpression(LOWEST))
	}

	return errorClasses
}
