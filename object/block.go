package object

import (
	"bytes"
	"emerald/ast"
	"strings"
)

type Block struct {
	*BaseEmeraldValue
	Parameters []ast.Expression
	Body       *ast.BlockStatement
	Env        Environment
}

func (b *Block) ParentClass() EmeraldValue { return nil }
func (*Block) Type() EmeraldValueType      { return BLOCK_VALUE }
func (b *Block) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range b.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("do ")
	out.WriteString("|")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString("|\n")
	out.WriteString(b.Body.String())
	out.WriteString("\nend")

	return out.String()
}

func NewBlock(params []ast.Expression, body *ast.BlockStatement, env Environment) *Block {
	return &Block{
		Parameters: params,
		Body:       body,
		Env:        env,
	}
}
