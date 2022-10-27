package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseCaseExpression() ast.Expression {
	caseExpr := &ast.CaseExpression{
		Token:       p.curToken,
		Subject:     nil,
		WhenClauses: []*ast.WhenClause{},
		Alternative: &ast.BlockStatement{
			Statements: []ast.Statement{},
		},
	}

	p.nextToken()

	caseExpr.Subject = p.parseExpression(LOWEST)
	p.nextToken()
	p.nextIfCurSemicolonOrNewline()

	for p.curTokenIs(lexer.WHEN) {
		whenClause := &ast.WhenClause{Matchers: []ast.Expression{}}

		whenClause.Matchers, _ = p.parseMethodArgsWithoutParentheses()

		whenClause.Consequence = p.parseBlockStatement(lexer.END, lexer.ELSE, lexer.WHEN)

		caseExpr.WhenClauses = append(caseExpr.WhenClauses, whenClause)
	}

	if p.curTokenIs(lexer.ELSE) {
		caseExpr.Alternative = p.parseBlockStatement(lexer.END)
	}

	return caseExpr
}
