package parser

import (
	"emerald/parser/ast"
)

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	switch left := left.(type) {
	case *ast.MethodCall:
		p.nextToken()

		return p.appendAssignmentToMethodCall(left, p.parseExpression(ASSIGN))
	default:
		node := &ast.AssignmentExpression{Token: p.curToken, Name: left}

		p.nextToken()

		node.Value = p.parseExpression(ASSIGN)

		return node
	}
}

func (p *Parser) appendAssignmentToMethodCall(left *ast.MethodCall, assignedExpression ast.Expression) ast.Expression {
	left.Method = ast.IdentifierExpression{
		Token: left.Method.Token,
		Value: left.Method.Value + "=",
	}
	left.Method.Token.Literal = left.Method.TokenLiteral() + "="

	left.Arguments = append(left.Arguments, assignedExpression)

	return left
}
