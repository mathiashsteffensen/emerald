package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type NextStatement struct {
	Token lexer.Token // the lexer.NEXT token
	Value Expression
}

func (ns *NextStatement) statementNode()       {}
func (ns *NextStatement) TokenLiteral() string { return ns.Token.Literal }
func (ns *NextStatement) String(indents ...int) string {
	var out strings.Builder

	indented(&out, indents[0], "next ")
	if ns.Value != nil {
		out.WriteString(ns.Value.String(0))
	}

	return out.String()
}
