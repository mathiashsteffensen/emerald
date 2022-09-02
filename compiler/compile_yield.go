package compiler

import "emerald/parser/ast"

func (c *Compiler) compileYield(node ast.Yield) error {
	for _, argument := range node.Arguments {
		err := c.Compile(argument)
		if err != nil {
			return err
		}
	}

	c.emit(OpYield, len(node.Arguments))

	return nil
}
