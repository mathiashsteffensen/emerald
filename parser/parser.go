package parser

import (
	"bytes"
	ast "emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l              *lexer.Lexer
	errors         []string
	curToken       lexer.Token
	peekToken      lexer.Token
	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.l.Run()

	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdentifierExpression)
	p.registerPrefix(lexer.INT, p.parseIntegerLiteral)
	p.registerPrefix(lexer.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(lexer.BANG, p.parsePrefixExpression)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.TRUE, p.parseBooleanExpression)
	p.registerPrefix(lexer.FALSE, p.parseBooleanExpression)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.IF, p.parseIfExpression)
	p.registerPrefix(lexer.UNLESS, p.parseUnlessExpression)
	p.registerPrefix(lexer.WHILE, p.parseWhileExpression)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.LBRACE, p.parseHashLiteral)
	p.registerPrefix(lexer.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(lexer.NULL, p.parseNullExpression)
	p.registerPrefix(lexer.DEF, p.parseMethodLiteral)
	p.registerPrefix(lexer.CLASS, p.parseClassLiteral)
	p.registerPrefix(lexer.COLON, p.parseSymbolLiteral)
	p.registerPrefix(lexer.INSTANCE_VAR, p.parseInstanceVariable)
	p.registerPrefix(lexer.GLOBAL_IDENT, p.parseGlobalVariable)
	p.registerPrefix(lexer.MODULE, p.parseModuleLiteral)
	p.registerPrefix(lexer.SELF, p.parseSelf)
	p.registerPrefix(lexer.REGEXP, p.parseRegexpLiteral)
	p.registerPrefix(lexer.YIELD, p.parseYield)
	p.registerPrefix(lexer.CASE, p.parseCaseExpression)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.ASTERISK, p.parseInfixExpression)
	p.registerInfix(lexer.ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(lexer.MATCH, p.parseInfixExpression)
	p.registerInfix(lexer.SPACESHIP, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.CASE_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LT_OR_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.GT_OR_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.BOOL_AND, p.parseInfixExpression)
	p.registerInfix(lexer.BOOL_OR, p.parseInfixExpression)
	p.registerInfix(lexer.BOOL_AND_ASSIGN, p.parseBoolAndAssign)
	p.registerInfix(lexer.BOOL_OR_ASSIGN, p.parseBoolOrAssign)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.DOT, p.parseMethodCall)
	p.registerInfix(lexer.LBRACKET, p.parseIndexAccessor)
	p.registerInfix(lexer.SCOPE, p.parseScopeAccessor)
	p.registerInfix(lexer.IF, p.parseIfModifier)
	p.registerInfix(lexer.UNLESS, p.parseUnlessModifier)
	p.registerInfix(lexer.WHILE, p.parseWhileModifier)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) unexpectedEofError() {
	p.addError("syntax error, unexpected end-of-input")
}

func (p *Parser) peekError(t lexer.TokenType) {
	p.addError(fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type))
}

func (p *Parser) peekErrorMultiple(types ...lexer.TokenType) {
	var typesBuffer bytes.Buffer

	for i, tokenType := range types {
		typesBuffer.WriteString(string(tokenType))

		if i != len(types)-1 {
			typesBuffer.WriteString(", ")
		}
	}

	msg := fmt.Sprintf("expected next token to be one of [%s], got %s instead", typesBuffer.String(), p.peekToken.Type)
	p.addError(msg)
}

func (p *Parser) noPrefixParseFnError() {
	msg := fmt.Sprintf(
		"no prefix parse function for %s found at line %d, column %d\n%s",
		p.curToken.Type,
		p.curToken.Line,
		p.curToken.Column,
		p.l.Snapshot(p.curToken),
	)
	p.addError(msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseAST() *ast.AST {
	program := &ast.AST{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		if p.curToken.Type != lexer.EOF {
			p.nextToken()
		}
	}

	return program
}

func (p *Parser) parseBlockStatement(endTokens ...lexer.TokenType) *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(endTokens...) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	p.expectCur(endTokens...)

	return block
}

func (p *Parser) parseStatement() ast.Statement {
	var result ast.Statement

	switch p.curToken.Type {
	case lexer.NEWLINE:
		result = nil
	case lexer.RETURN:
		result = p.parseReturnStatement()
	case lexer.BREAK:
		result = p.parseBreakStatement()
	default:
		result = p.parseExpressionStatement()
	}

	p.nextIfSemicolonOrNewline()

	return result
}

func (p *Parser) parseBoolModifierFromStatement(stmt ast.Statement) ast.Expression {
	switch p.curToken.Type {
	case lexer.IF:
		return p.parseIfModifierFromStatement(stmt)
	default:
		return nil
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if p.curToken.Type == lexer.EOF {
		return nil
	}

	leftExp := p.parseAsPrefix()

	for !p.peekTokenIs(lexer.SEMICOLON) && !p.peekTokenIs(lexer.NEWLINE) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseAsPrefix() ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError()
		return nil
	}

	return prefix()
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer", p.curToken.Literal))
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as float", p.curToken.Literal))
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()

	p.nextToken()

	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBooleanExpression() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(lexer.TRUE)}
}

func (p *Parser) parseNullExpression() ast.Expression {
	return &ast.NullExpression{Token: p.curToken}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseExpressionList(delim lexer.TokenType) []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(delim) {
		p.nextToken()
		return args
	}

	p.nextToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(delim) {
		return nil
	}

	return args
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{
		Token: p.curToken,
	}

	value := p.parseExpressionList(lexer.RBRACKET)

	arr.Value = value

	return arr
}

func (p *Parser) parseStaticClassLiteral() ast.Expression {
	class := &ast.StaticClassLiteral{Token: p.curToken}

	p.nextToken()

	if !p.expectPeek(lexer.SELF) {
		return nil
	}

	p.nextIfSemicolonOrNewline()

	class.Body = p.parseBlockStatement(lexer.END)

	return class
}

func (p *Parser) curTokenIs(ts ...lexer.TokenType) bool {
	for _, t := range ts {
		if p.curToken.Type == t {
			return true
		}
	}
	return false
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekTokenIsMultiple(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.peekTokenIs(t) {
			return true
		}
	}

	return false
}

func (p *Parser) expectCur(ts ...lexer.TokenType) bool {
	if p.curTokenIs(ts...) {
		return true
	} else {
		p.peekError(ts[0])
		return false
	}
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) expectPeekMultiple(advance bool, types ...lexer.TokenType) bool {
	if p.peekTokenIsMultiple(types...) {
		if advance {
			p.nextToken()
		}
		return true
	}

	p.peekErrorMultiple(types...)
	return false
}

func (p *Parser) nextIfSemicolonOrNewline() {
	if p.peekTokenIs(lexer.SEMICOLON) || p.peekTokenIs(lexer.NEWLINE) {
		p.nextToken()
	}
}

func (p *Parser) nextIfCurSemicolonOrNewline() {
	if p.curTokenIs(lexer.SEMICOLON) || p.curTokenIs(lexer.NEWLINE) {
		p.nextToken()
	}
}

func (p *Parser) nextIfNewline() {
	if p.peekTokenIs(lexer.NEWLINE) {
		p.nextToken()
	}
}
