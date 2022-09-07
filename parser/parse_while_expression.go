package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: p.curToken}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	p.nextIfSemicolonOrNewline()

	expression.Consequence = &ast.BlockStatement{Token: p.curToken}
	expression.Consequence.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			expression.Consequence.Statements = append(expression.Consequence.Statements, stmt)
		}

		p.nextToken()
	}

	p.nextIfSemicolonOrNewline()

	return expression
}

func (p *Parser) parseWhileModifier(consequence ast.Expression) ast.Expression {
	expression := &ast.WhileExpression{Token: p.curToken, Consequence: &ast.BlockStatement{
		Statements: []ast.Statement{&ast.ExpressionStatement{
			Expression: consequence,
		}},
	}}

	p.nextToken()

	expression.Condition = p.parseExpression(MODIFIER)

	return expression
}
