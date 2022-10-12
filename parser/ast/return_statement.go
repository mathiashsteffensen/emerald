package ast

import (
	"emerald/parser/lexer"
	"strings"
)

type ReturnStatement struct {
	Token       lexer.Token // the lexer.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String(indents ...int) string {
	var out strings.Builder

	indented(&out, indents[0], "return ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String(0))
	}

	return out.String()
}
