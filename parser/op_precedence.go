package parser

import (
	"emerald/parser/lexer"
)

const (
	_ int = iota
	LOWEST
	MODIFIER   // val = 10 if true
	BOOL_OR    // ||
	BOOL_AND   // &&
	COMPARATOR // ==
	ORDERING   // > or <
	SUM        // +
	PRODUCT    // *
	PREFIX     // -X or !X
	CALL       // myFunction(X)
	ACCESSOR   // myHash.property
)

var precedences = map[lexer.TokenType]int{
	lexer.IF:        MODIFIER,
	lexer.BOOL_OR:   BOOL_OR,
	lexer.BOOL_AND:  BOOL_AND,
	lexer.MATCH:     COMPARATOR,
	lexer.SPACESHIP: COMPARATOR,
	lexer.EQ:        COMPARATOR,
	lexer.NOT_EQ:    COMPARATOR,
	lexer.LT:        ORDERING,
	lexer.GT:        ORDERING,
	lexer.LT_OR_EQ:  ORDERING,
	lexer.GT_OR_EQ:  ORDERING,
	lexer.PLUS:      SUM,
	lexer.MINUS:     SUM,
	lexer.SLASH:     PRODUCT,
	lexer.ASTERISK:  PRODUCT,
	lexer.LPAREN:    CALL,
	lexer.LBRACKET:  CALL,
	lexer.DOT:       ACCESSOR,
	lexer.SCOPE:     ACCESSOR,
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
