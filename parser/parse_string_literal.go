package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
	"regexp"
)

func (p *Parser) parseStringLiteral() ast.Expression {
	str := p.newString()

	if p.peekTokenIs(lexer.LTEMPLATE) {
		strTemplate := &ast.StringTemplate{Chain: &ast.StringTemplateChainString{
			StringLiteral: str,
			Next:          p.parseStringTemplateExpression(),
			First:         true,
		}}

		return strTemplate
	}

	return str
}

func (p *Parser) parseStringLiteralTemplateRecursive() *ast.StringTemplateChainString {
	str := p.newString()

	strTemplate := &ast.StringTemplateChainString{
		StringLiteral: str,
		First:         false,
	}

	if p.peekTokenIs(lexer.LTEMPLATE) {
		strTemplate.Next = p.parseStringTemplateExpression()
	}

	return strTemplate
}

func (p *Parser) parseStringTemplateExpression() *ast.StringTemplateChainExpression {
	p.nextToken()
	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(lexer.RBRACE) {
		return nil
	}

	templateExpression := &ast.StringTemplateChainExpression{
		Expression: expr,
		Next:       nil,
	}

	if p.peekTokenIs(lexer.STRING) {
		p.nextToken()
		templateExpression.Next = p.parseStringLiteralTemplateRecursive()
	}

	return templateExpression
}

type stringEscapeFunc func(str string) string

func newStringEscapeFunc(escapeToken string, raw string) stringEscapeFunc {
	re := regexp.MustCompile(fmt.Sprintf(`\\%s`, escapeToken))

	return func(str string) string {
		return re.ReplaceAllLiteralString(str, raw)
	}
}

var stringEscapeFunctions = []stringEscapeFunc{
	newStringEscapeFunc("n", "\n"),
	newStringEscapeFunc("t", "\t"),
	newStringEscapeFunc("a", "\a"),
	newStringEscapeFunc("b", "\b"),
	newStringEscapeFunc("v", "\v"),
	newStringEscapeFunc("f", "\f"),
	newStringEscapeFunc("r", "\r"),
	newStringEscapeFunc("s", " "),
}

func escapeString(str string) string {
	for _, function := range stringEscapeFunctions {
		str = function(str)
	}

	return str
}

func (p *Parser) newString() *ast.StringLiteral {
	return &ast.StringLiteral{Token: p.curToken, Value: escapeString(p.curToken.Literal)}
}
