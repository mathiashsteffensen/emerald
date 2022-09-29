package parser

import "emerald/parser/ast"

func (p *Parser) parseBoolOrAssign(left ast.Expression) ast.Expression {
	precedence := p.curPrecedence()

	p.nextToken()

	right := p.parseExpression(precedence)

	var alternative ast.Expression

	switch left := left.(type) {
	case *ast.MethodCall:
		alternative = p.appendAssignmentToMethodCall(left.Dup(), right)
	default:
		alternative = &ast.AssignmentExpression{
			Name:  left,
			Value: right,
		}
	}

	return &ast.IfExpression{
		Condition:   left,
		Consequence: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: left}}},
		Alternative: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: alternative}}},
	}
}

func (p *Parser) parseBoolAndAssign(left ast.Expression) ast.Expression {
	precedence := p.curPrecedence()

	p.nextToken()

	right := p.parseExpression(precedence)

	var consequence ast.Expression

	switch left := left.(type) {
	case *ast.MethodCall:
		consequence = p.appendAssignmentToMethodCall(left.Dup(), right)
	default:
		consequence = &ast.AssignmentExpression{
			Name:  left,
			Value: right,
		}
	}

	return &ast.IfExpression{
		Condition:   left,
		Consequence: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: consequence}}},
		Alternative: &ast.BlockStatement{Statements: []ast.Statement{&ast.ExpressionStatement{Expression: left}}},
	}
}
