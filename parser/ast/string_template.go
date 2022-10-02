package ast

import "bytes"

type StringTemplateChainString struct {
	*StringLiteral
	Next  *StringTemplateChainExpression
	First bool
}

func (s *StringTemplateChainString) String() string {
	str := s.StringLiteral.String()

	if s.Next == nil {
		return str[1:]
	}

	var startIndex int

	if s.First {
		startIndex = 0
	} else {
		startIndex = 1
	}

	return str[startIndex:len(str)-1] + s.Next.String()
}

type StringTemplateChainExpression struct {
	Expression
	Next *StringTemplateChainString
}

func (s *StringTemplateChainExpression) String() string {
	var out bytes.Buffer

	out.WriteString("#{")
	out.WriteString(s.Expression.String())
	out.WriteString("}")

	if s.Next == nil {
		out.Write([]byte{'"'})
	} else {
		out.WriteString(s.Next.String())
	}

	return out.String()
}

type StringTemplate struct {
	Chain *StringTemplateChainString
}

func (s *StringTemplate) expressionNode() {}

func (s *StringTemplate) TokenLiteral() string {
	return s.Chain.TokenLiteral()
}

func (s *StringTemplate) String() string {
	return s.Chain.String()
}
