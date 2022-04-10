package evaluator

import (
	"emerald/ast"
	"emerald/object"
	"fmt"
)

func Eval(executionContext object.ExecutionContext, node ast.Node, env object.Environment) object.EmeraldValue {
	switch node := node.(type) {
	case *ast.AST:
		return evalAST(executionContext, node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(executionContext, node, env)
	case *ast.ExpressionStatement:
		return Eval(executionContext, node.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(executionContext, node.ReturnValue, env)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	case *ast.AssignmentExpression:
		val := Eval(executionContext, node.Value, env)

		if isError(val) {
			return val
		}

		env.Set(node.Name.Value, val)

		return val
	case *ast.ClassLiteral:
		return evalClassLiteral(executionContext, node, env)
	case *ast.StaticClassLiteral:
		return evalClassLiteral(executionContext, node, env)
	case *ast.MethodLiteral:
		return executionContext.Target.DefineMethod(
			executionContext.IsStatic,
			object.NewBlock(node.Parameters, node.Body, env),
			object.NewString(node.Name.String()),
		)
	case *ast.IntegerLiteral:
		return object.NewInteger(node.Value)
	case *ast.StringLiteral:
		return object.NewString(node.Value)
	case *ast.SymbolLiteral:
		return evalSymbolLiteral(node)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(executionContext, node, env)
	case *ast.BlockLiteral:
		return object.NewBlock(node.Parameters, node.Body, object.NewEnclosedEnvironment(env))
	case *ast.BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(executionContext, node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(executionContext, node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(executionContext, node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right, executionContext)
	case *ast.IfExpression:
		return evalIfExpression(executionContext, node, env)
	case *ast.IdentifierExpression:
		ident := evalIdentifier(executionContext, node, env)

		switch ident.(type) {
		case *object.Block:
			return evalBlock(executionContext, executionContext.Target, node.Value, ident, nil, []object.EmeraldValue{})
		case *object.WrappedBuiltInMethod:
			return evalBlock(executionContext, executionContext.Target, node.Value, ident, nil, []object.EmeraldValue{})
		}

		return ident
	case *ast.MethodCall:
		target := Eval(executionContext, node.Left, env)
		if isError(target) {
			return target
		}

		return evalCallExpression(executionContext, target, node.CallExpression, env)
	case *ast.CallExpression:
		return evalCallExpression(executionContext, executionContext.Target, node, env)
	case *ast.NullExpression:
		return object.NULL
	default:
		return newError("Unimplemented ast expression %T", node)
	}
}

func Yield(executionContext object.ExecutionContext) object.YieldFunc {
	return func(block *object.Block, args ...object.EmeraldValue) object.EmeraldValue {
		return evalBlock(executionContext, executionContext.Target, "YIELD", block, nil, args)
	}
}

func evalAST(executionContext object.ExecutionContext, program *ast.AST, env object.Environment) object.EmeraldValue {
	var result object.EmeraldValue

	for _, statement := range program.Statements {
		result = Eval(executionContext, statement, env)

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

func evalBlockStatement(executionContext object.ExecutionContext, block *ast.BlockStatement, env object.Environment) object.EmeraldValue {
	var result object.EmeraldValue

	for _, statement := range block.Statements {
		result = Eval(executionContext, statement, env)
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
	ctx object.ExecutionContext,
) object.EmeraldValue {
	return left.SEND(Eval, Yield(ctx), operator, left, nil, right)
}

func evalIfExpression(
	executionContext object.ExecutionContext,
	ie *ast.IfExpression,
	env object.Environment,
) object.EmeraldValue {
	condition := Eval(executionContext, ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(executionContext, ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(executionContext, ie.Alternative, env)
	} else {
		return object.NULL
	}
}

func evalIdentifier(
	executionContext object.ExecutionContext,
	node *ast.IdentifierExpression,
	env object.Environment,
) object.EmeraldValue {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}

	method, err := executionContext.Target.ExtractMethod(node.Value, executionContext.Target, executionContext.Target)
	if err != nil {
		fmt.Printf("%#v\n", env)
		return err
	}

	return method
}

func evalExpressions(
	executionContext object.ExecutionContext,
	exps []ast.Expression,
	env object.Environment,
) []object.EmeraldValue {
	var result []object.EmeraldValue

	for _, e := range exps {
		evaluated := Eval(executionContext, e, env)
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

func evalBlock(executionContext object.ExecutionContext, target object.EmeraldValue, name string, block object.EmeraldValue, givenBlock *object.Block, args []object.EmeraldValue) object.EmeraldValue {
	switch block := block.(type) {
	case *object.Block:
		extendedEnv := object.ExtendBlockEnv(block.Env, block.Parameters, args)
		evaluated := Eval(object.ExecutionContext{Target: target}, block.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.WrappedBuiltInMethod:
		return block.Method(target, givenBlock, Yield(executionContext), args...)
	default:
		return newError("not a method: %s %T", name, block)
	}
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
