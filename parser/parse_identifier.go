package parser

import (
	"emerald/ast"
	"emerald/lexer"
)

func (p *Parser) parseIdentifierExpression() ast.Expression {
	node := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}

	return p.parseIdentifierOrAssignment(node)
}

func (p *Parser) parseInstanceVariable() ast.Expression {
	node := &ast.InstanceVariable{IdentifierExpression: &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}}

	return p.parseIdentifierOrAssignment(node)
}

func (p *Parser) parseGlobalVariable() ast.Expression {
	node := &ast.GlobalVariable{IdentifierExpression: &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}}

	return p.parseIdentifierOrAssignment(node)
}

func (p *Parser) parseIdentifierOrAssignment(identNode ast.Expression) ast.Expression {
	switch p.peekToken.Type {
	case lexer.ASSIGN:
		p.nextToken()

		return p.parseAssignmentExpression(identNode)
	case lexer.BOOL_AND_ASSIGN, lexer.BOOL_OR_ASSIGN:
		p.nextToken()

		infix := &ast.InfixExpression{
			Token:    p.curToken,
			Operator: p.curToken.Literal[:len(p.curToken.Literal)-1],
			Left:     identNode,
		}

		infix.Right = p.parseAssignmentExpression(identNode)

		return infix
	default:
		return identNode
	}
}
