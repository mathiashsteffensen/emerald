package parser

import (
	"emerald/parser/ast"
	"fmt"
)

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	switch left := left.(type) {
	case *ast.MethodCall:
		p.nextToken()

		return p.appendAssignmentToMethodCall(left, p.parseExpression(ASSIGN))
	default:
		p.checkAssignmentTarget(left)
		p.declareLocalIfIdentifier(left)

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
	left.Assignment = true

	return left
}

func (p *Parser) parseInfixOperatorAssignment(infixOp string) infixParseFn {
	return func(left ast.Expression) ast.Expression {
		p.checkAssignmentTarget(left)

		p.declareLocalIfIdentifier(left)

		node := &ast.AssignmentExpression{Token: p.curToken, Name: left}

		p.nextToken()

		node.Value = &ast.InfixExpression{
			Token:    node.Token,
			Left:     left,
			Operator: infixOp,
			Right:    p.parseExpression(ASSIGN),
		}

		return node
	}
}

// checkAssignmentTarget reports a target that the compiler cannot store to,
// such as the array in `foo [1] = 2`. The current token is the operator.
func (p *Parser) checkAssignmentTarget(target ast.Expression) {
	switch target.(type) {
	case ast.IdentifierExpression, *ast.InstanceVariable, *ast.GlobalVariable, *ast.MethodCall:
		return
	}

	msg := fmt.Sprintf(
		"syntax error, unexpected '%s' at line %d, column %d\n%s",
		p.curToken.Literal,
		p.curToken.Line,
		p.curToken.Column,
		p.l.Snapshot(p.curToken),
	)
	p.addError(msg)
}
