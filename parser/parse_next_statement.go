package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseNextStatement() ast.Statement {
	stmt := &ast.NextStatement{Token: p.curToken, Value: &ast.NullExpression{Token: p.curToken}}
	if !p.peekTokenIsMultiple(lexer.NEWLINE, lexer.SEMICOLON, lexer.END, lexer.RBRACE, lexer.RESCUE, lexer.EOF, lexer.IF, lexer.UNLESS) {
		p.nextToken()
		stmt.Value = p.parseExpression(MODIFIER)
	}

	if !p.curTokenIs(lexer.NEWLINE) && p.peekTokenIsMultiple(lexer.IF, lexer.UNLESS) {
		p.nextToken()
		return &ast.ExpressionStatement{Token: stmt.Token, Expression: p.parseBoolModifierFromStatement(stmt)}
	}

	return stmt
}
