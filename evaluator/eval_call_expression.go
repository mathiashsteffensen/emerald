package evaluator

import (
	"emerald/ast"
	"emerald/object"
	"fmt"
)

func evalCallExpression(executionContext object.ExecutionContext, target object.EmeraldValue, node *ast.CallExpression, env object.Environment) object.EmeraldValue {
	method, err := target.ExtractMethod(node.Method.(*ast.IdentifierExpression).Value, target, target)
	if err != nil {
		fmt.Printf("%#v\n", env)
		return err
	}

	args := evalExpressions(executionContext, node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	var block *object.Block

	if node.Block == nil {
		block = nil
	} else {
		evaluated := Eval(executionContext, node.Block, env)
		if isError(evaluated) {
			return evaluated
		}

		block = evaluated.(*object.Block)
	}

	return evalBlock(executionContext, target, node.Method.String(), method, block, args)
}
