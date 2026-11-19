package parser

import (
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
)

func (p *Parser) parseIdentifierExpression() ast.Expression {
	node := ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(lexer.LBRACE) || p.peekTokenStartsDoBlock() {
		callExpression := ast.CallExpression{
			Token:     p.curToken,
			Method:    node,
			Arguments: []ast.Expression{},
		}

		callExpression.Block = p.parseBlockLiteral()

		return callExpression
	}

	if p.peekTokenIs(lexer.ASSIGN) {
		p.nextToken()
		return p.parseAssignmentExpression(node)
	}

	if p.peekTokenDoesntSignifyCallArguments() || p.peekTokenIndexesLocal(node.Value) {
		return node
	}

	callExpression := ast.CallExpression{
		Token:  p.curToken,
		Method: node,
	}

	callExpression.Arguments, callExpression.KeywordArguments = p.parseMethodArgsWithoutParentheses()
	callExpression.Block = p.parseBlockLiteral()

	return callExpression
}

// peekTokenIndexesLocal is true when `[` after a local is an index even after
// a space: a local cannot take command arguments.
func (p *Parser) peekTokenIndexesLocal(name string) bool {
	return p.peekTokenIs(lexer.LBRACKET) && p.isLocalVariable(name)
}

func (p *Parser) parseInstanceVariable() ast.Expression {
	return &ast.InstanceVariable{IdentifierExpression: ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}}
}

func (p *Parser) parseGlobalVariable() ast.Expression {
	return &ast.GlobalVariable{IdentifierExpression: ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}}
}

func (p *Parser) parseSelf() ast.Expression {
	return &ast.Self{IdentifierExpression: ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}}
}
