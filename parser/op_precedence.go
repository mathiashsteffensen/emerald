package parser

import (
	"emerald/lexer"
)

const (
	_ int = iota
	LOWEST
	MODIFIER    // val = 10 if true
	BOOL_OR     // ||
	BOOL_AND    // &&
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
	lexer.BOOL_OR:  BOOL_OR,
	lexer.BOOL_AND: BOOL_AND,
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
	lexer.LBRACKET: CALL,
	lexer.DOT:      ACCESSOR,
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
