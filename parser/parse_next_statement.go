package parser

import "emerald/parser/ast"

func (p *Parser) parseNextStatement() ast.Statement {
	stmt := &ast.NextStatement{Token: p.curToken, Value: &ast.NullExpression{Token: p.curToken}}
	if !p.peekTokenEndsJumpValue() {
		p.nextToken()
		stmt.Value = p.parseExpression(MODIFIER)
	}

	return p.parseStatementModifiers(stmt.Token, stmt)
}
