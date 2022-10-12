package compiler

import (
	ast "emerald/parser/ast"
	"sort"
)

func (c *Compiler) compileHashLiteral(node *ast.HashLiteral) error {
	keys := []ast.Expression{}
	for k := range node.Value {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String(0) < keys[j].String(0)
	})
	for _, k := range keys {
		err := c.Compile(k)
		if err != nil {
			return err
		}
		err = c.Compile(node.Value[k])
		if err != nil {
			return err
		}
	}
	c.emit(OpHash, len(node.Value)*2)

	return nil
}
