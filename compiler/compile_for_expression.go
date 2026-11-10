package compiler

import (
	"emerald/parser/ast"
)

// compileForExpression compiles `for x in items; body; end` as
// `items.each { |element| x = element; body }`. The block shares the locals
// of the enclosing scope, and `break` and `next` work as in any block.
func (c *Compiler) compileForExpression(node *ast.ForExpression) {
	// `for` is a keyword, so the program cannot name this parameter
	element := &ast.IdentifierExpression{Token: node.Token, Value: "for"}

	assignVariable := &ast.ExpressionStatement{
		Token: node.Variable.Token,
		Expression: &ast.AssignmentExpression{
			Token: node.Variable.Token,
			Name:  node.Variable,
			Value: *element,
		},
	}

	body := &ast.BlockStatement{
		Token:      node.Body.Token,
		Statements: append([]ast.Statement{assignVariable}, node.Body.Statements...),
	}

	c.compileMethodCall(&ast.MethodCall{
		Token: node.Token,
		Left:  node.Iterable,
		CallExpression: ast.CallExpression{
			Token:     node.Token,
			Method:    ast.IdentifierExpression{Token: node.Token, Value: "each"},
			Arguments: []ast.Expression{},
			Block: &ast.BlockLiteral{
				Token:       node.Token,
				Arguments:   []*ast.IdentifierExpression{element},
				Body:        body,
				ForLoopBody: true,
			},
		},
	})
}
