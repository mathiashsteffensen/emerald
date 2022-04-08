package evaluator

import (
	"emerald/ast"
	"emerald/object"
)

func evalArrayLiteral(
	executionContext object.ExecutionContext,
	al *ast.ArrayLiteral,
	env object.Environment,
) object.EmeraldValue {
	values := evalExpressions(executionContext, al.Value, env)

	return object.NewArray(values)
}
