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

	for !p.curTokenIs(lexer.END, lexer.ELSE, lexer.ELSIF, lexer.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			expression.Consequence.Statements = append(expression.Consequence.Statements, stmt)
		}

		p.nextToken()
	}

	for p.curTokenIs(lexer.ELSIF) {
		elseIf := ast.ElseIf{Consequence: &ast.BlockStatement{Token: p.curToken}}

		p.nextToken()

		elseIf.Condition = p.parseExpression(LOWEST)

		p.nextIfSemicolonOrNewline()
		p.nextToken()

		for !p.curTokenIs(lexer.END, lexer.ELSE, lexer.ELSIF, lexer.EOF) {
			stmt := p.parseStatement()

			if stmt != nil {
				elseIf.Consequence.Statements = append(elseIf.Consequence.Statements, stmt)
			}

			p.nextToken()
		}

		expression.ElseIfs = append(expression.ElseIfs, elseIf)
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

func (p *Parser) parseUnlessExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	p.nextIfSemicolonOrNewline()

	expression.Alternative = &ast.BlockStatement{Token: p.curToken}
	expression.Alternative.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(lexer.END, lexer.ELSE, lexer.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			expression.Alternative.Statements = append(expression.Alternative.Statements, stmt)
		}

		p.nextToken()
	}

	if p.curTokenIs(lexer.ELSE) {
		p.nextIfSemicolonOrNewline()

		expression.Consequence = p.parseBlockStatement(lexer.END)
	}

	p.nextIfSemicolonOrNewline()

	return expression
}

func (p *Parser) parseUnlessModifier(alternative ast.Expression) ast.Expression {
	return p.parseUnlessModifierFromStatement(&ast.ExpressionStatement{
		Expression: alternative,
	})
}

func (p *Parser) parseUnlessModifierFromStatement(stmt ast.Statement) ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken, Alternative: &ast.BlockStatement{Statements: []ast.Statement{stmt}}}

	p.nextToken()

	expression.Condition = p.parseExpression(MODIFIER)

	return expression
}

func (p *Parser) parseTernary(condition ast.Expression) ast.Expression {
	expression := &ast.IfExpression{
		Token:       p.curToken,
		Condition:   condition,
		Consequence: &ast.BlockStatement{},
		Alternative: &ast.BlockStatement{},
	}

	p.nextToken()

	consequence := &ast.ExpressionStatement{Expression: p.parseExpression(LOWEST)}
	expression.Consequence.Statements = append(expression.Consequence.Statements, consequence)

	if !p.expectPeek(lexer.COLON) {
		return nil
	}

	p.nextToken()

	alternative := &ast.ExpressionStatement{Expression: p.parseExpression(LOWEST)}
	expression.Alternative.Statements = append(expression.Alternative.Statements, alternative)

	return expression
}
