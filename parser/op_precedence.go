package parser

import (
	"emerald/lexer"
)

const (
	_ int = iota
	LOWEST
	MODIFIER    // val = 10 if true
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	ACCESSOR    // myHash.property
)

var precedences = map[lexer.TokenType]int{
	lexer.IF:       MODIFIER,
	lexer.EQ:       EQUALS,
	lexer.NOT_EQ:   EQUALS,
	lexer.LT:       LESSGREATER,
	lexer.GT:       LESSGREATER,
	lexer.LT_OR_EQ: LESSGREATER,
	lexer.GT_OR_EQ: LESSGREATER,
	lexer.PLUS:     SUM,
	lexer.MINUS:    SUM,
	lexer.SLASH:    PRODUCT,
	lexer.ASTERISK: PRODUCT,
	lexer.LPAREN:   CALL,
	lexer.DOT:      ACCESSOR,
	lexer.LBRACKET: ACCESSOR,
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
