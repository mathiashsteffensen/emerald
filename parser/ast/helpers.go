package ast

import "strings"

func indented(out *strings.Builder, indent int, str string) *strings.Builder {
	out.WriteString(strings.Repeat("	", indent))
	out.WriteString(str)
	return out
}
