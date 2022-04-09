package parser

import (
	"bytes"
	"emerald/ast"
	"emerald/lexer"
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
	p.registerPrefix(lexer.BANG, p.parsePrefixExpression)
	p.registerPrefix(lexer.MINUS, p.parsePrefixExpression)
	p.registerPrefix(lexer.TRUE, p.parseBooleanExpression)
	p.registerPrefix(lexer.FALSE, p.parseBooleanExpression)
	p.registerPrefix(lexer.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(lexer.IF, p.parseIfExpression)
	p.registerPrefix(lexer.STRING, p.parseStringLiteral)
	p.registerPrefix(lexer.LBRACE, p.parseHashLiteral)
	p.registerPrefix(lexer.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(lexer.NULL, p.parseNullExpression)
	p.registerPrefix(lexer.DEF, p.parseMethodLiteral)
	p.registerPrefix(lexer.CLASS, p.parseClassLiteral)
	p.registerPrefix(lexer.COLON, p.parseSymbolLiteral)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.PLUS, p.parseInfixExpression)
	p.registerInfix(lexer.MINUS, p.parseInfixExpression)
	p.registerInfix(lexer.SLASH, p.parseInfixExpression)
	p.registerInfix(lexer.ASTERISK, p.parseInfixExpression)
	p.registerInfix(lexer.EQ, p.parseInfixExpression)
	p.registerInfix(lexer.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LT, p.parseInfixExpression)
	p.registerInfix(lexer.GT, p.parseInfixExpression)
	p.registerInfix(lexer.LT_OR_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.GT_OR_EQ, p.parseInfixExpression)
	p.registerInfix(lexer.LPAREN, p.parseCallExpression)
	p.registerInfix(lexer.DOT, p.parseMethodCall)

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

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.addError(msg)
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

func (p *Parser) noPrefixParseFnError(t lexer.TokenType) {
	msg := fmt.Sprintf(
		"no prefix parse function for %s found at line %d, column %d\n%s",
		t,
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

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(lexer.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifierExpression() ast.Expression {
	node := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(lexer.ASSIGN) {
		p.nextToken()

		return p.parseAssignmentExpression(node)
	}

	return node
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseSymbolLiteral() ast.Expression {
	symbol := &ast.SymbolLiteral{Token: p.curToken}

	if !p.expectPeekMultiple(lexer.IDENT, lexer.STRING) {
		return nil
	}

	symbol.Value = p.curToken.Literal

	return symbol
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
	return &ast.BooleanExpression{Token: p.curToken, Value: p.curTokenIs(lexer.TRUE)}
}

func (p *Parser) parseNullExpression() ast.Expression {
	return &ast.NullExpression{Token: p.curToken}
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	node := &ast.AssignmentExpression{}

	if left, ok := left.(*ast.IdentifierExpression); !ok {
		p.addError(fmt.Sprintf("cannot assign to expression of type %T", left))
		return nil
	} else {
		node.Name = left
		node.Token = left.Token
	}

	p.nextToken()

	node.Value = p.parseExpression(LOWEST)

	return node
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()

	expression.Condition = p.parseExpression(MODIFIER)

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	expression.Consequence = &ast.BlockStatement{Token: p.curToken}
	expression.Consequence.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(lexer.END) && !p.curTokenIs(lexer.ELSE) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			expression.Consequence.Statements = append(expression.Consequence.Statements, stmt)
		}

		p.nextToken()
	}

	if p.curTokenIs(lexer.ELSE) {
		if p.peekTokenIs(lexer.SEMICOLON) {
			p.nextToken()
		}

		expression.Alternative = p.parseBlockStatement(lexer.END)
	}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	return expression
}

func (p *Parser) parseBlockStatement(endToken lexer.TokenType) *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(endToken) && !p.curTokenIs(lexer.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseCallExpression(method ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Method: method}

	exp.Arguments = p.parseCallArguments()

	exp.Block = p.parseBlockLiteral()

	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	return p.parseExpressionList(lexer.RPAREN)
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

func (p *Parser) parseHashLiteral() ast.Expression {
	value := make(map[string]ast.Expression)

	for p.curToken.Type != lexer.RBRACE {
		if p.peekTokenIs(lexer.RBRACE) {
			p.nextToken()
			break
		}

		if !p.expectPeekMultiple(lexer.IDENT, lexer.STRING, lexer.INT) {
			return nil
		}

		key := p.curToken.Literal

		if !p.expectPeek(lexer.COLON) {
			return nil
		}

		p.nextToken()

		value[key] = p.parseExpression(LOWEST)

		if !p.peekTokenIs(lexer.COMMA) {
			if !p.expectPeek(lexer.RBRACE) {
				return nil
			}
		} else {
			p.nextToken()
		}
	}

	return &ast.HashLiteral{
		Value: value,
		Token: p.curToken,
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arr := &ast.ArrayLiteral{
		Token: p.curToken,
	}

	value := p.parseExpressionList(lexer.RBRACKET)

	arr.Value = value

	return arr
}

func (p *Parser) parseMethodLiteral() ast.Expression {
	method := &ast.MethodLiteral{Token: p.curToken}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	method.Name = p.parseIdentifierExpression()

	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken()

		method.Parameters = p.parseCallArguments()
	} else {
		method.Parameters = make([]ast.Expression, 0)
	}

	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	method.Body = p.parseBlockStatement(lexer.END)

	return method
}

func (p *Parser) parseClassLiteral() ast.Expression {
	if p.peekTokenIs(lexer.APPEND) {
		return p.parseStaticClassLiteral()
	}

	class := &ast.ClassLiteral{Token: p.curToken}

	p.nextToken()

	class.Name = p.parseIdentifierExpression().(*ast.IdentifierExpression)
	if p.peekTokenIs(lexer.SEMICOLON) {
		p.nextToken()
	}

	class.Body = p.parseBlockStatement(lexer.END)

	return class
}

func (p *Parser) parseStaticClassLiteral() ast.Expression {
	class := &ast.StaticClassLiteral{Token: p.curToken}

	p.nextToken()

	if !p.expectPeek(lexer.SELF) {
		return nil
	}

	class.Body = p.parseBlockStatement(lexer.END)

	return class
}

func (p *Parser) parseMethodCall(left ast.Expression) ast.Expression {
	node := &ast.MethodCall{Token: p.curToken, Left: left, CallExpression: &ast.CallExpression{}}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}

	node.Method = p.parseIdentifierExpression()

	if p.peekTokenIs(lexer.LPAREN) {
		p.nextToken()
		node.Arguments = p.parseCallArguments()
	}

	node.Block = p.parseBlockLiteral()

	if p.curTokenIs(lexer.DOT) {
		return p.parseMethodCall(node)
	}

	return node
}

func (p *Parser) parseBlockLiteral() *ast.BlockLiteral {
	if !p.peekTokenIs(lexer.LBRACE) && !p.peekTokenIs(lexer.DO) {
		return nil
	}

	var endToken lexer.TokenType
	if p.peekTokenIs(lexer.LBRACE) {
		endToken = lexer.RBRACE
	} else {
		endToken = lexer.END
	}

	p.nextToken()

	block := &ast.BlockLiteral{Body: &ast.BlockStatement{}, Token: p.curToken}

	if p.peekTokenIs(lexer.LINE) {
		p.nextToken()
		block.Parameters = p.parseExpressionList(lexer.LINE)
	}

	block.Body = p.parseBlockStatement(endToken)

	return block
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
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

func (p *Parser) expectPeekMultiple(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.peekTokenIs(t) {
			p.nextToken()
			return true
		}
	}

	p.peekErrorMultiple(types...)
	return false
}
