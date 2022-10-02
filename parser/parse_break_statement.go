package parser

import (
	"emerald/parser/ast"
)

func (p *Parser) parseBreakStatement() ast.Statement {
	node := &ast.BreakStatement{
		Token: p.curToken,
	}

	p.nextToken()

	if p.curPrecedence() == MODIFIER {
		return &ast.ExpressionStatement{Expression: p.parseBoolModifierFromStatement(node)}
	}

	node.Value = p.parseAsPrefix()
	p.nextToken()

	if p.curPrecedence() == MODIFIER {
		return &ast.ExpressionStatement{Expression: p.parseBoolModifierFromStatement(node)}
	}

	return node
}
