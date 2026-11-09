package parser

import (
	"emerald/parser/ast"
	"fmt"
)

func (p *Parser) parseCallExpression(method ast.Expression) ast.Expression {
	methodName, isMethodName := method.(ast.IdentifierExpression)
	if !isMethodName {
		p.notCallableError()
		return nil
	}

	exp := ast.CallExpression{Token: p.curToken, Method: methodName}

	exp.Arguments, exp.KeywordArguments = p.parseMethodArgsWithParentheses()

	exp.Block = p.parseBlockLiteral()

	return exp
}

func (p *Parser) notCallableError() {
	msg := fmt.Sprintf(
		"syntax error, unexpected '(' at line %d, column %d\n%s",
		p.curToken.Line,
		p.curToken.Column,
		p.l.Snapshot(p.curToken),
	)
	p.addError(msg)
}
