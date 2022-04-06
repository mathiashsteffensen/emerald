package evaluator

import (
	"emerald/ast"
	"emerald/object"
)

func evalClassLiteral(
	cl *ast.ClassLiteral,
	env object.Environment,
) object.EmeraldValue {
	name := cl.Name.Value

	class, ok := env.Get(name)
	if !ok {
		class = object.NewClass(name, object.Object, object.BuiltInMethodSet{})
		env.Set(name, class)
	}

	return Eval(class, cl.Body, env)
}
