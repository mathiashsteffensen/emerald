package evaluator

import (
	"emerald/ast"
	"emerald/object"
)

func evalSymbolLiteral(sl *ast.SymbolLiteral) object.EmeraldValue {
	return object.NewSymbol(sl.Value)
}
