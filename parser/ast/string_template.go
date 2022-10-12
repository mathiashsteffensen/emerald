package ast

import (
	"reflect"
	"strings"
)

type StringTemplateChainString struct {
	*StringLiteral
	Next  *StringTemplateChainExpression
	First bool
}

func (s *StringTemplateChainString) next() Next {
	return s.Next
}

func (s *StringTemplateChainString) String(indents ...int) string {
	str := s.StringLiteral.String(0)

	if s.Next == nil {
		return str[1:]
	}

	var startIndex int

	if s.First {
		startIndex = 0
	} else {
		startIndex = 1
	}

	return str[startIndex:len(str)-1] + s.Next.String(0)
}

type StringTemplateChainExpression struct {
	Expression
	Next *StringTemplateChainString
}

func (s *StringTemplateChainExpression) next() Next {
	return s.Next
}

func (s *StringTemplateChainExpression) String(indents ...int) string {
	var out strings.Builder

	out.WriteString("#{")
	out.WriteString(s.Expression.String(0))
	out.WriteString("}")

	if s.Next == nil {
		out.Write([]byte{'"'})
	} else {
		out.WriteString(s.Next.String(0))
	}

	return out.String()
}

type StringTemplate struct {
	Chain *StringTemplateChainString
}

type Next interface {
	next() Next
}

func (s *StringTemplate) Count() int {
	var next Next = s.Chain

	count := 0

	for !reflect.ValueOf(next).IsNil() {
		count += 1
		next = next.next()
	}

	return count
}

func (s *StringTemplate) expressionNode() {}

func (s *StringTemplate) TokenLiteral() string {
	return s.Chain.TokenLiteral()
}

func (s *StringTemplate) String(indents ...int) string {
	return s.Chain.String(0)
}
