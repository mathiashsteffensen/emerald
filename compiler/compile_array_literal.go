package compiler

import "emerald/ast"

func (c *Compiler) compileArrayLiteral(node *ast.ArrayLiteral) error {
	for _, val := range node.Value {
		err := c.Compile(val)
		if err != nil {
			return err
		}
	}

	c.emit(OpArray, len(node.Value))

	return nil
}
