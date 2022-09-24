package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	p.nextIfSemicolonOrNewline()

	expression.Consequence = &ast.BlockStatement{Token: p.curToken}
	expression.Consequence.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.ELSE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			expression.Consequence.Statements = append(expression.Consequence.Statements, stmt)
		}

		p.nextToken()
	}

	if p.curTokenIs(lexer.ELSE) {
		p.nextIfSemicolonOrNewline()

		expression.Alternative = p.parseBlockStatement(lexer.END)
	}

	p.nextIfSemicolonOrNewline()

	return expression
}

func (p *Parser) parseIfModifier(consequence ast.Expression) ast.Expression {
	return p.parseIfModifierFromStatement(&ast.ExpressionStatement{
		Expression: consequence,
	})
}

func (p *Parser) parseIfModifierFromStatement(stmt ast.Statement) ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken, Consequence: &ast.BlockStatement{Statements: []ast.Statement{stmt}}}

	p.nextToken()

	expression.Condition = p.parseExpression(MODIFIER)

	return expression
}
