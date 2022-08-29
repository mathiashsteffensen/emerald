package ast

import (
	"bytes"
	"emerald/parser/lexer"
)

type AssignmentExpression struct {
	Token lexer.Token // the lexer.IDENT token
	Name  Expression
	Value Expression
}

func (ls *AssignmentExpression) expressionNode()      {}
func (ls *AssignmentExpression) TokenLiteral() string { return ls.Token.Literal }
func (ls *AssignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	out.WriteString(ls.Value.String())

	return out.String()
}
