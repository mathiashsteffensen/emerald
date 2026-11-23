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

	expression.Consequence = p.parseBlockStatement(lexer.END)

	return expression
}

func (p *Parser) parseWhileModifier(consequence ast.Expression) ast.Expression {
	return p.parseWhileModifierFromStatement(&ast.ExpressionStatement{
		Expression: consequence,
	})
}

func (p *Parser) parseWhileModifierFromStatement(stmt ast.Statement) ast.Expression {
	expression := &ast.WhileExpression{Token: p.curToken, Consequence: &ast.BlockStatement{Statements: []ast.Statement{stmt}}}

	p.nextToken()

	expression.Condition = p.parseExpression(MODIFIER)

	return expression
}
