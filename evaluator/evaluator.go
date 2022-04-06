package evaluator

import (
	"emerald/ast"
	"emerald/object"
	"fmt"
)

func Eval(definitionContext object.EmeraldValue, node ast.Node, env object.Environment) object.EmeraldValue {
	switch node := node.(type) {
	case *ast.AST:
		return evalAST(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(definitionContext, node, env)
	case *ast.ExpressionStatement:
		return Eval(definitionContext, node.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(definitionContext, node.ReturnValue, env)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	case *ast.AssignmentExpression:
		val := Eval(definitionContext, node.Value, env)

		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)

		return val
	case *ast.MethodLiteral:
		return definitionContext.DefineMethod(object.NewBlock(node.Parameters, node.Body, env), object.NewString(node.Name.String()))
	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)
	case *ast.StringLiteral:
		return object.NewString(node.Value)
	case *ast.BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(definitionContext, node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(definitionContext, node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(definitionContext, node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right, env)
	case *ast.IfExpression:
		return evalIfExpression(definitionContext, node, env)
	case *ast.IdentifierExpression:
		return evalIdentifier(node, env)
	case *ast.CallExpression:
		function := Eval(definitionContext, node.Method, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(definitionContext, node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return evalBlock(definitionContext, node.Method.String(), function, args)
	case *ast.NullExpression:
		return object.NULL
	}

	return nil
}

func evalAST(program *ast.AST, env object.Environment) object.EmeraldValue {
	var result object.EmeraldValue

	for _, statement := range program.Statements {
		result = Eval(object.Object, statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		default:
			if isError(result) {
				return result
			}
		}
	}
	return result
}

func evalBlockStatement(definitionContext object.EmeraldValue, block *ast.BlockStatement, env object.Environment) object.EmeraldValue {
	var result object.EmeraldValue

	for _, statement := range block.Statements {
		result = Eval(definitionContext, statement, env)
		if result != nil {
			if result.Type() == object.RETURN_VALUE || isError(result) {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.EmeraldValue) object.EmeraldValue {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Inspect())
	}
}

func evalBangOperatorExpression(right object.EmeraldValue) object.EmeraldValue {
	return nativeBoolToBooleanObject(!isTruthy(right))
}

func evalMinusPrefixOperatorExpression(right object.EmeraldValue) object.EmeraldValue {
	if right.ParentClass().(*object.Class).Name != "Integer" {
		return newError("unknown operator: -%s", right.ParentClass().(*object.Class).Name)
	}

	value := right.(*object.IntegerInstance).Value
	return object.NewInteger(-value)
}

func evalInfixExpression(
	operator string,
	left, right object.EmeraldValue,
	env object.Environment,
) object.EmeraldValue {
	return left.SEND(Eval, env, operator, left, nil, right)
}

func evalIfExpression(definitionContext object.EmeraldValue, ie *ast.IfExpression, env object.Environment) object.EmeraldValue {
	condition := Eval(definitionContext, ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(definitionContext, ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(definitionContext, ie.Alternative, env)
	} else {
		return object.NULL
	}
}

func evalIdentifier(
	node *ast.IdentifierExpression,
	env object.Environment,
) object.EmeraldValue {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}

	if !object.Object.RespondsTo(node.Value, object.Object) {
		return newError("identifier not found: " + node.Value)
	}

	method, _ := object.Object.ExtractMethod(node.Value, object.Object)

	return method
}

func evalExpressions(
	definitionContext object.EmeraldValue,
	exps []ast.Expression,
	env object.Environment,
) []object.EmeraldValue {
	var result []object.EmeraldValue

	for _, e := range exps {
		evaluated := Eval(definitionContext, e, env)
		if isError(evaluated) {
			return []object.EmeraldValue{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func newError(format string, a ...interface{}) object.EmeraldValue {
	return object.NewStandardError(fmt.Sprintf(format, a...))
}

func isTruthy(obj object.EmeraldValue) bool {
	switch obj {
	case object.NULL:
		return false
	case object.TRUE:
		return true
	case object.FALSE:
		return false
	default:
		return true
	}
}

func evalBlock(definitionContext object.EmeraldValue, name string, block object.EmeraldValue, args []object.EmeraldValue) object.EmeraldValue {
	switch block := block.(type) {
	case *object.Block:
		extendedEnv := extendBlockEnv(block, args)
		evaluated := Eval(definitionContext, block.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.WrappedBuiltInMethod:
		return block.Method(object.Object, nil, args...)
	default:
		return newError("not a method: %s", name)
	}
}

func extendBlockEnv(
	fn *object.Block,
	args []object.EmeraldValue,
) object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.(*ast.IdentifierExpression).Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.EmeraldValue) object.EmeraldValue {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func isError(obj object.EmeraldValue) bool {
	if obj != nil {
		return obj.ParentClass() == object.StandardError
	}
	return false
}

func nativeBoolToBooleanObject(input bool) object.EmeraldValue {
	if input {
		return object.TRUE
	}
	return object.FALSE
}
