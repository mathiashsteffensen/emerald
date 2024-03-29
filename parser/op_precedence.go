package parser

import (
	"emerald/parser/lexer"
)

// https://docs.ruby-lang.org/en/3.1/syntax/precedence_rdoc.html

const (
	_ int = iota
	LOWEST
	MODIFIER   // val = 10 if true
	ASSIGN     // =
	TERNARY    // ?, :
	BOOL_OR    // ||
	BOOL_AND   // &&
	COMPARATOR // ==
	ORDERING   // > or <
	BIN_SHIFT  // << or >>
	SUM        // + or -
	PRODUCT    // *
	PREFIX     // -X or !X
	CALL       // myFunction(X)
	ACCESSOR   // myHash.property
)

var precedences = map[lexer.TokenType]int{
	lexer.WHILE:           MODIFIER,
	lexer.IF:              MODIFIER,
	lexer.UNLESS:          MODIFIER,
	lexer.BOOL_OR_ASSIGN:  ASSIGN,
	lexer.BOOL_AND_ASSIGN: ASSIGN,
	lexer.PLUS_ASSIGN:     ASSIGN,
	lexer.MINUS_ASSIGN:    ASSIGN,
	lexer.SLASH_ASSIGN:    ASSIGN,
	lexer.ASTERISK_ASSIGN: ASSIGN,
	lexer.ASSIGN:          ASSIGN,
	lexer.QUESTION:        TERNARY,
	lexer.BOOL_OR:         BOOL_OR,
	lexer.BOOL_AND:        BOOL_AND,
	lexer.MATCH:           COMPARATOR,
	lexer.SPACESHIP:       COMPARATOR,
	lexer.EQ:              COMPARATOR,
	lexer.CASE_EQ:         COMPARATOR,
	lexer.NOT_EQ:          COMPARATOR,
	lexer.LT:              ORDERING,
	lexer.GT:              ORDERING,
	lexer.LT_OR_EQ:        ORDERING,
	lexer.GT_OR_EQ:        ORDERING,
	lexer.APPEND:          BIN_SHIFT,
	lexer.PLUS:            SUM,
	lexer.MINUS:           SUM,
	lexer.SLASH:           PRODUCT,
	lexer.ASTERISK:        PRODUCT,
	lexer.LPAREN:          CALL,
	lexer.LBRACKET:        CALL,
	lexer.DOT:             ACCESSOR,
	lexer.SCOPE:           ACCESSOR,
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
