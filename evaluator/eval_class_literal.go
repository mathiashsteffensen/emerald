package evaluator

import (
	"emerald/ast"
	"emerald/object"
)

func evalClassLiteral(
	executionContext object.ExecutionContext,
	cl ast.Expression,
	env object.Environment,
) object.EmeraldValue {
	switch cl := cl.(type) {
	case *ast.ClassLiteral:
		name := cl.Name.Value

		class, ok := env.Get(name)
		if !ok {
			class = object.NewClass(name, object.Object, object.BuiltInMethodSet{}, object.BuiltInMethodSet{})
			env.Set(name, class)
		}

		return Eval(object.ExecutionContext{Target: class}, cl.Body, env)
	case *ast.StaticClassLiteral:
		return Eval(object.ExecutionContext{Target: executionContext.Target, IsStatic: true}, cl.Body, env)
	default:
		return newError("invalid class passed to evalClassLiteral, got=%T", cl)
	}
}
