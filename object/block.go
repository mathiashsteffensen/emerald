package object

import (
	"emerald/ast"
	"fmt"
)

type Block struct {
	*BaseEmeraldValue
	Parameters   []ast.Expression
	Body         *ast.BlockStatement
	Instructions []byte
	NumLocals    int
}

func (b *Block) ParentClass() EmeraldValue { return nil }
func (*Block) Type() EmeraldValueType      { return BLOCK_VALUE }
func (b *Block) Inspect() string           { return fmt.Sprintf("#<Block:%p>", b) }

func NewBlock(params []ast.Expression, instructions []byte, numLocals int) *Block {
	return &Block{
		Parameters:   params,
		Instructions: instructions,
		NumLocals:    numLocals,
	}
}

func ExtendBlockEnv(
	env Environment,
	params []ast.Expression,
	args []EmeraldValue,
) Environment {
	env = NewEnclosedEnvironment(env)
	for paramIdx, param := range params {
		env.Set(param.(*ast.IdentifierExpression).Value, args[paramIdx])
	}
	return env
}
