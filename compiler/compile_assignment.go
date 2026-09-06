package compiler

import (
	"emerald/bytecode"
	"emerald/heap"
	ast "emerald/parser/ast"
	"unicode"
)

func (c *Compiler) compileAssignment(node *ast.AssignmentExpression) {
	if target, ok := node.Name.(*ast.MethodCall); ok {
		value := node.Value.(*ast.InfixExpression)
		c.Compile(target.Left)
		c.emit(bytecode.OpNull, node.Token)
		for _, arg := range target.Arguments {
			c.Compile(arg)
		}
		c.emit(bytecode.OpDupN, node.Token, len(target.Arguments)+2)
		c.emit(bytecode.OpSend, node.Token, c.addConstant(c.rt.NewSymbol(target.Method.Value)), len(target.Arguments), 0)
		skip := -1
		if value.Operator == "||" || value.Operator == "&&" {
			c.emit(bytecode.OpDupN, node.Token, 1)
			jump := bytecode.OpJumpTruthy
			if value.Operator == "&&" {
				jump = bytecode.OpJumpNotTruthy
			}
			skip = c.emit(jump, node.Token, 9999)
			c.emit(bytecode.OpPop, node.Token)
			c.emit(bytecode.OpPop, node.Token)
			c.Compile(value.Right)
		} else {
			c.compileCallExpression(ast.CallExpression{
				Token: node.Token, Method: ast.IdentifierExpression{Value: value.Operator}, Arguments: []ast.Expression{value.Right},
			})
		}
		c.emit(bytecode.OpSendAssign, node.Token, c.addConstant(c.rt.NewSymbol(target.Method.Value+"=")), len(target.Arguments)+1, 0)
		if skip != -1 {
			end := c.emit(bytecode.OpJump, node.Token, 9999)
			c.changeOperand(skip, len(c.currentInstructions()))
			c.emit(bytecode.OpDropN, node.Token, len(target.Arguments)+2)
			c.changeOperand(end, len(c.currentInstructions()))
		}
		return
	}

	c.Compile(node.Value)

	switch name := node.Name.(type) {
	case ast.IdentifierExpression:
		if unicode.IsUpper(rune(name.Value[0])) {
			c.emit(bytecode.OpConstantSet, node.Token, c.addConstant(c.rt.NewSymbol(name.Value)))
			return
		}

		stringName := name.String(0)

		symbol, ok := c.symbolTable.Resolve(stringName)
		if !ok {
			symbol = c.symbolTable.Define(stringName)
		}

		switch symbol.Scope {
		case heap.GlobalScope:
			c.emit(bytecode.OpSetGlobal, node.Token, symbol.Index)
		case heap.LocalScope:
			c.emit(bytecode.OpSetLocal, node.Token, symbol.Index)
		case heap.FreeScope:
			c.emit(bytecode.OpSetFree, node.Token, symbol.Index)
		}
	case *ast.InstanceVariable:
		c.emit(bytecode.OpInstanceVarSet, node.Token, c.addConstant(c.rt.NewSymbol(name.Value)))
	case *ast.GlobalVariable:
		symbol, ok := c.symbolTable.Resolve(name.Value)
		if !ok {
			symbol = c.symbolTable.DefineGlobal(name.Value)
		}

		c.emit(bytecode.OpSetGlobal, node.Token, symbol.Index)
	}
}
