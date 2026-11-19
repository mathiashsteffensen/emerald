package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
)

// localScope holds the local variable names known at one point of the source,
// so that the parser can tell `x [0]` (an index) from `puts [0]` (a call).
type localScope struct {
	names map[string]bool

	// visibleEnclosing is the scope around a block. It is nil for the top
	// level, `def`, `class` and `module`, which cannot see outer locals.
	visibleEnclosing *localScope
}

// enterIsolatedLocalScope starts a scope that sees no enclosing locals. Call
// the returned func to go back to the enclosing scope.
func (p *Parser) enterIsolatedLocalScope() (leave func()) {
	return p.replaceLocalScope(&localScope{names: map[string]bool{}})
}

// enterBlockLocalScope starts a scope that also sees the enclosing locals. Call
// the returned func to go back to the enclosing scope.
func (p *Parser) enterBlockLocalScope() (leave func()) {
	return p.replaceLocalScope(&localScope{names: map[string]bool{}, visibleEnclosing: p.locals})
}

func (p *Parser) replaceLocalScope(scope *localScope) (restore func()) {
	enclosing := p.locals
	p.locals = scope

	return func() { p.locals = enclosing }
}

func (p *Parser) declareLocal(name string) {
	p.locals.names[name] = true
}

// declareLocalIfIdentifier declares an assignment target. Other targets, such
// as constants, instance variables and method calls, are not locals.
func (p *Parser) declareLocalIfIdentifier(target ast.Expression) {
	identifier, isIdentifier := target.(ast.IdentifierExpression)

	if isIdentifier && !lexer.IsConstantName(identifier.Value) {
		p.declareLocal(identifier.Value)
	}
}

// declareParameters declares parameters that have no default value, such as
// keywords. Positional parameters are declared before their default is parsed.
func (p *Parser) declareParameters(parameters []*ast.IdentifierExpression) {
	for _, parameter := range parameters {
		p.declareLocal(parameter.Value)
	}
}

// DeclareLocalVariables declares locals that an earlier parse defined at the
// top level, as in a REPL, so that `x [0]` there still indexes x.
func (p *Parser) DeclareLocalVariables(names []string) {
	for _, name := range names {
		p.declareLocal(name)
	}
}

func (p *Parser) isLocalVariable(name string) bool {
	for scope := p.locals; scope != nil; scope = scope.visibleEnclosing {
		if scope.names[name] {
			return true
		}
	}

	return false
}
