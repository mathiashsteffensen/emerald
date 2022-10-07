package ast

import "bytes"

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type AST struct {
	Statements []Statement
}

func (ast *AST) TokenLiteral() string {
	if len(ast.Statements) > 0 {
		return ast.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (ast *AST) String() string {
	var out bytes.Buffer
	for _, s := range ast.Statements {
		out.WriteString(s.String())
		if len(ast.Statements) != 1 {
			out.WriteString("\n")
		}
	}
	return out.String()
}
