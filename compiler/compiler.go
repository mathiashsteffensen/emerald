package compiler

import (
	"context"
	"emerald/bytecode"
	"emerald/core"
	"emerald/heap"
	"emerald/object"
	"emerald/parser"
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"fmt"
	"reflect"
)

type EmittedInstruction struct {
	Opcode   bytecode.Opcode
	Position int
}

type Compiler struct {
	ctx          context.Context
	err          error
	rt           *core.Runtime
	instructions bytecode.Instructions
	opCount      int
	symbolTable  *heap.SymbolTable
	scopes       []CompilationScope
	scopeIndex   int
	aborted      bool
}

type ConstructorOption func(c *Compiler)

func New(l *lexer.Lexer, rt *core.Runtime, options ...ConstructorOption) *Compiler {
	return NewContext(context.Background(), l, rt, options...)
}

func NewContext(ctx context.Context, l *lexer.Lexer, rt *core.Runtime, options ...ConstructorOption) *Compiler {
	if ctx == nil {
		ctx = context.Background()
	}

	mainScope := CompilationScope{
		bytecode: bytecode.Bytecode{
			Instructions: bytecode.Instructions{},
			DebugTokens:  map[int]lexer.Token{},
			Lexer:        l,
		},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c := &Compiler{
		ctx:          ctx,
		rt:           rt,
		instructions: bytecode.Instructions{},
		symbolTable:  rt.Heap.SymbolTable,
		scopes:       []CompilationScope{mainScope},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *Compiler) SetLexer(l *lexer.Lexer) {
	c.scopes[c.scopeIndex].bytecode.Lexer = l
}

func (c *Compiler) SetRuntime(rt *core.Runtime) {
	c.rt = rt
	c.symbolTable = rt.Heap.SymbolTable
}

func Compile(fileName string, content string, rt *core.Runtime) *bytecode.Bytecode {
	bc, _ := CompileContext(context.Background(), fileName, content, rt)

	return bc
}

func CompileContext(ctx context.Context, fileName string, content string, rt *core.Runtime) (*bytecode.Bytecode, error) {
	bc, _, err := compileContext(ctx, fileName, content, rt)
	return bc, err
}

func compile(fileName string, content string, rt *core.Runtime) (*bytecode.Bytecode, bool) {
	bc, ok, _ := compileContext(context.Background(), fileName, content, rt)
	return bc, ok
}

func compileContext(ctx context.Context, fileName string, content string, rt *core.Runtime) (*bytecode.Bytecode, bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	l := lexer.NewContext(ctx, lexer.NewInput(fileName, content))
	defer l.Close()

	p := parser.NewContext(ctx, l)
	ast := p.ParseAST()

	if err := p.Err(); err != nil {
		return nil, false, err
	}

	if len(p.Errors()) != 0 {
		rt.Raise(rt.NewException(fmt.Sprintf("failed to parse source file %s\n\n%s", fileName, p.Errors()[0])))
		return emptyBytecode(l), false, nil
	}

	c := NewContext(ctx, l, rt)
	c.Compile(ast)

	if err := c.Err(); err != nil {
		return nil, false, err
	}

	if c.aborted {
		return emptyBytecode(l), false, nil
	}

	bc := c.Bytecode()

	return bc, true, nil
}

func CompileBlock(fileName string, content string, rt *core.Runtime) *bytecode.Bytecode {
	bc, ok := compile(fileName, content, rt)
	if !ok {
		return bc
	}

	bc.Instructions = append(bc.Instructions, byte(bytecode.OpReturn))

	return bc
}

func emptyBytecode(l *lexer.Lexer) *bytecode.Bytecode {
	return &bytecode.Bytecode{
		Instructions: bytecode.Instructions{},
		DebugTokens:  map[int]lexer.Token{},
		Lexer:        l,
	}
}

func (c *Compiler) Err() error {
	c.stopIfCanceled()
	return c.err
}

func (c *Compiler) stopIfCanceled() bool {
	if c.aborted {
		return true
	}

	if err := c.ctx.Err(); err != nil {
		c.err = err
		c.aborted = true
		c.clearBytecode()
		return true
	}

	return false
}

func (c *Compiler) clearBytecode() {
	for i := range c.scopes {
		c.scopes[i].bytecode.Instructions = bytecode.Instructions{}
		c.scopes[i].bytecode.DebugTokens = map[int]lexer.Token{}
		c.scopes[i].lastInstruction = EmittedInstruction{}
		c.scopes[i].previousInstruction = EmittedInstruction{}
	}
}

func (c *Compiler) abort(message string) {
	if c.aborted {
		return
	}

	c.aborted = true
	c.rt.Raise(c.rt.NewException(message))
	c.clearBytecode()
}

func (c *Compiler) Compile(node ast.Node) {
	if c.stopIfCanceled() {
		return
	}

	if isNilNode(node) {
		c.abort("compiler received nil AST node")
		return
	}

	switch node := node.(type) {
	case *ast.AST:
		for _, s := range node.Statements {
			c.Compile(s)
			if c.aborted {
				return
			}
		}
	case *ast.ExpressionStatement:
		c.Compile(node.Expression)
		if c.aborted {
			return
		}
		c.emit(bytecode.OpPop, node.Token)
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			c.Compile(s)
			if c.aborted {
				return
			}
		}
	case *ast.ReturnStatement:
		c.Compile(node.ReturnValue)
		if c.aborted {
			return
		}
		c.emit(bytecode.OpReturnValue, node.Token)
	case *ast.PrefixExpression:
		c.compilePrefixExpression(node)
	case *ast.AssignmentExpression:
		c.compileAssignment(node)
	case *ast.Self:
		c.emit(bytecode.OpSelf, node.Token)
	case ast.Yield:
		c.compileYield(node)
	case ast.IdentifierExpression:
		c.compileIdentifierExpression(node)
	case *ast.InstanceVariable:
		c.compileIdentifierExpression(node)
	case *ast.GlobalVariable:
		c.compileIdentifierExpression(node)
	case ast.CallExpression:
		c.emit(bytecode.OpSelf, node.Token) // Method calls without a receiver has an implicit self receiver
		c.compileCallExpression(node)
	case *ast.MethodCall:
		c.compileMethodCall(node)
	case *ast.ScopeAccessor:
		c.compileScopeAccessor(node)
	case *ast.InfixExpression:
		c.compileInfixExpression(node)
	case *ast.IfExpression:
		c.compileIfExpression(node)
	case *ast.CaseExpression:
		c.compileCaseExpression(node)

		if c.lastInstructionIs(bytecode.OpPop) {
			c.removeLastPop()
		}
	case *ast.WhileExpression:
		c.compileWhileExpression(node)
	case *ast.IntegerLiteral:
		integer := c.rt.NewInteger(node.Value)
		c.emit(bytecode.OpPushConstant, node.Token, c.addConstant(integer))
	case *ast.FloatLiteral:
		float := c.rt.NewFloat(node.Value)
		c.emit(bytecode.OpPushConstant, node.Token, c.addConstant(float))
	case *ast.BooleanLiteral:
		if node.Value {
			c.emit(bytecode.OpTrue, node.Token)
		} else {
			c.emit(bytecode.OpFalse, node.Token)
		}
	case *ast.NullExpression:
		c.emit(bytecode.OpNull, node.Token)
	case *ast.StringLiteral:
		c.compileStringLiteral(node)
	case *ast.StringTemplate:
		c.compileStringTemplate(node)
	case *ast.SymbolLiteral:
		sym := c.rt.NewSymbol(node.Value)
		c.emit(bytecode.OpPushConstant, node.Token, c.addConstant(sym))
	case *ast.RegexpLiteral:
		c.compileRegexpLiteral(node)
	case *ast.ArrayLiteral:
		c.compileArrayLiteral(node)
	case *ast.HashLiteral:
		c.compileHashLiteral(node)
	case *ast.MethodLiteral:
		c.compileMethodLiteral(node)
	case *ast.ClassLiteral:
		c.compileClassLiteral(node)
	case *ast.StaticClassLiteral:
		c.compileStaticClassLiteral(node)
	case *ast.ModuleLiteral:
		c.compileModuleLiteral(node)
	}
}

func isNilNode(node ast.Node) bool {
	if node == nil {
		return true
	}

	value := reflect.ValueOf(node)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func (c *Compiler) compileStatementsWithReturnValue(statements []ast.Statement, debugToken lexer.Token) {
	if len(statements) == 0 {
		c.emit(bytecode.OpNull, debugToken)
	} else {
		for _, s := range statements {
			c.Compile(s)
		}
	}
}

func (c *Compiler) Bytecode() *bytecode.Bytecode {
	return &bytecode.Bytecode{
		Instructions: c.currentInstructions(),
		DebugTokens:  c.currentScope().bytecode.DebugTokens,
		Lexer:        c.currentScope().bytecode.Lexer,
	}
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		bytecode: bytecode.Bytecode{
			Instructions: bytecode.Instructions{},
			DebugTokens:  map[int]lexer.Token{},
			Lexer:        c.currentScope().bytecode.Lexer,
		},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = heap.NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() bytecode.Bytecode {
	bytecode := c.currentScope().bytecode

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer

	return bytecode
}

func (c *Compiler) lastInstructionIs(op bytecode.Opcode) bool {
	if c.aborted {
		return false
	}

	if len(c.currentInstructions()) == 0 {
		return false
	}
	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	if c.aborted {
		return
	}

	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction
	old := c.currentInstructions()

	c.scopes[c.scopeIndex].bytecode.Instructions = old[:last.Position]
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	if c.aborted {
		return
	}

	ins := c.currentInstructions()
	if pos < 0 || pos+len(newInstruction) > len(ins) {
		c.abort("compiler attempted to replace an invalid instruction range")
		return
	}

	for i := 0; i < len(newInstruction); i++ {
		ins[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) replaceLastInstructionWith(op bytecode.Opcode) {
	if c.aborted {
		return
	}

	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, bytecode.Make(op))
	c.scopes[c.scopeIndex].lastInstruction.Opcode = op
}

func (c *Compiler) changeOperand(opPos int, operands ...int) {
	if c.aborted {
		return
	}

	if opPos < 0 || opPos >= len(c.currentInstructions()) {
		c.abort("compiler attempted to change an invalid instruction operand")
		return
	}

	op := bytecode.Opcode(c.currentInstructions()[opPos])
	newInstruction := bytecode.Make(op, operands...)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) emit(op bytecode.Opcode, debugToken lexer.Token, operands ...int) int {
	if c.stopIfCanceled() {
		return -1
	}

	ins := bytecode.Make(op, operands...)
	pos := c.addInstruction(ins, debugToken)

	c.setLastInstruction(op, pos)

	c.opCount += 1

	return pos
}

func (c *Compiler) emitConstantGet(name string, debugToken lexer.Token) {
	symbol := c.rt.NewSymbol(name)

	c.emit(bytecode.OpConstantGet, debugToken, c.addConstant(symbol))
}

func (c *Compiler) emitSymbol(symbol heap.Symbol, debugToken lexer.Token) {
	switch symbol.Scope {
	case heap.GlobalScope:
		c.emit(bytecode.OpGetGlobal, debugToken, symbol.Index)
	case heap.LocalScope:
		c.emit(bytecode.OpGetLocal, debugToken, symbol.Index)
	case heap.FreeScope:
		c.emit(bytecode.OpGetFree, debugToken, symbol.Index)
	}
}

// returns the instructions for the current CompilationScope
func (c *Compiler) currentInstructions() bytecode.Instructions {
	return c.currentScope().bytecode.Instructions
}

func (c *Compiler) currentScope() CompilationScope {
	return c.scopes[c.scopeIndex]
}

// addInstruction adds instructions to the instruction stack and returns its location
func (c *Compiler) addInstruction(ins []byte, debugToken lexer.Token) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].bytecode.Instructions = updatedInstructions
	c.scopes[c.scopeIndex].bytecode.DebugTokens[posNewInstruction] = debugToken

	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op bytecode.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

// addConstant adds a constant to the constant stack and returns its location
func (c *Compiler) addConstant(obj object.EmeraldValue) int {
	if c.stopIfCanceled() {
		return -1
	}

	return c.rt.Heap.AddConstant(obj)
}
