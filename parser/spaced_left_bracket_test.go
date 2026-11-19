package parser

import (
	"emerald/parser/ast"
	"emerald/parser/lexer"
	"testing"
)

func TestSpacedLeftBracketStartsCommandArgument(t *testing.T) {
	program := testParseAST(t, "puts [1, 2].to_s")

	expectStatementLength(t, program.Statements, 1)

	testExpressionStatement(t, program.Statements[0], func(call ast.CallExpression) {
		if len(call.Arguments) != 1 {
			t.Fatalf("got %d arguments, want 1", len(call.Arguments))
		}

		argument := call.Arguments[0].(*ast.MethodCall)
		if argument.Method.Value != "to_s" {
			t.Fatalf("argument method = %q, want to_s", argument.Method.Value)
		}

		testArrayLiteral(t, argument.Left.(*ast.ArrayLiteral), []any{1, 2})
	})
}

func TestUnspacedLeftBracketIndexesCallResult(t *testing.T) {
	for _, name := range []string{"foo", "café"} {
		program := testParseAST(t, name+"[1]")

		testExpressionStatement(t, program.Statements[0], func(expression *ast.MethodCall) {
			testMethodCall(t, expression, name, "[]", []any{1}, []string{}, false)
		})
	}
}

func TestSpacedLeftBracketAfterMethodNameStartsArgument(t *testing.T) {
	program := testParseAST(t, "obj.foo [1]\nobj.foo[1]")

	expectStatementLength(t, program.Statements, 2)

	testExpressionStatement(t, program.Statements[0], func(call *ast.MethodCall) {
		if call.Method.Value != "foo" || len(call.Arguments) != 1 {
			t.Fatalf("got method %q with %d arguments", call.Method.Value, len(call.Arguments))
		}

		testArrayLiteral(t, call.Arguments[0].(*ast.ArrayLiteral), []any{1})
	})

	testExpressionStatement(t, program.Statements[1], func(index *ast.MethodCall) {
		if index.Method.Value != "[]" {
			t.Fatalf("method = %q, want []", index.Method.Value)
		}

		testIntegerLiteral(t, index.Arguments[0], 1)
	})
}

func TestSpacedLeftBracketAfterYieldStartsArgument(t *testing.T) {
	program := testParseAST(t, "yield [1]")

	testExpressionStatement(t, program.Statements[0], func(yield ast.Yield) {
		if len(yield.Arguments) != 1 {
			t.Fatalf("got %d arguments, want 1", len(yield.Arguments))
		}

		testArrayLiteral(t, yield.Arguments[0].(*ast.ArrayLiteral), []any{1})
	})
}

func TestSpacedLeftBracketAfterWhenStartsMatcher(t *testing.T) {
	program := testParseAST(t, "case 1\nwhen [1]\n2\nend")

	testExpressionStatement(t, program.Statements[0], func(expression *ast.CaseExpression) {
		matchers := expression.WhenClauses[0].Matchers
		if len(matchers) != 1 {
			t.Fatalf("got %d matchers, want 1", len(matchers))
		}

		testArrayLiteral(t, matchers[0].(*ast.ArrayLiteral), []any{1})
	})
}

func TestNewlineBeforeLeftBracketStartsNewStatement(t *testing.T) {
	program := testParseAST(t, "foo\n[1]")

	expectStatementLength(t, program.Statements, 2)
}

// Each input must parse to the same tree as its unambiguous spelling.
func TestSpacedLeftBracketDependsOnLocalVariables(t *testing.T) {
	tests := []struct{ name, input, sameAs string }{
		{"assignment", "x = [1]; x [0]", "x = [1]; x[0]"},
		{"own assignment value", "x = x [0]", "x = x[0]"},
		{"operator assignment", "x += 1; x [0]", "x += 1; x[0]"},
		{"or assignment", "x ||= 1; x [0]", "x ||= 1; x[0]"},
		{"and assignment", "x &&= 1; x [0]", "x &&= 1; x[0]"},
		{"method parameter", "def m(a); a [0]; end", "def m(a); a[0]; end"},
		{"method parameter without parentheses", "def m a\na [0]\nend", "def m a\na[0]\nend"},
		{"earlier parameter in a default", "def m(a, b = a [0]); end", "def m(a, b = a[0]); end"},
		{"keyword parameter", "def m(a:); a [0]; end", "def m(a:); a[0]; end"},
		{"keyword parameter without parentheses", "def m a:\na [0]\nend", "def m a:\na[0]\nend"},
		{"method block parameter", "def m(&a); a [0]; end", "def m(&a); a[0]; end"},
		{"block parameter", "m { |a| a [0] }", "m { |a| a[0] }"},
		{"block keyword parameter", "m { |a:| a [0] }", "m { |a:| a[0] }"},
		{"do block parameter", "m do |a| a [0] end", "m do |a| a[0] end"},
		{"enclosing local in a block", "x = 1; m { x [0] }", "x = 1; m { x[0] }"},
		{"for variable", "for i in y; i [0]; end", "for i in y; i[0]; end"},
		{"for variable after the loop", "for i in y; end; i [0]", "for i in y; end; i[0]"},
		{"rescue variable", "def m; 1; rescue => e; e [0]; end", "def m; 1; rescue => e; e[0]; end"},
		{"non-ASCII local", "café = [1]; café [0]", "café = [1]; café[0]"},
		{"parameter in its own default", "def m(a = a [0]); end", "def m(a = a[0]); end"},

		{"method name", "foo [1]", "foo([1])"},
		{"non-ASCII method name", "café [1]", "café([1])"},
		{"constant", "Foo [1]", "Foo([1])"},
		{"assigned constant", "Foo = [1]; Foo [1]", "Foo = [1]; Foo([1])"},
		{"local of another method", "def a; x = 1; end; def b; x [0]; end", "def a; x = 1; end; def b; x([0]); end"},
		{"top-level local in a method", "x = 1; def m; x [0]; end", "x = 1; def m; x([0]); end"},
		{"top-level local in a class", "x = 1; class C; x [0]; end", "x = 1; class C; x([0]); end"},
		{"top-level local in a singleton class", "x = 1; class << self; x [0]; end", "x = 1; class << self; x([0]); end"},
		{"top-level local in a module", "x = 1; module M; x [0]; end", "x = 1; module M; x([0]); end"},
		{"block local after the block", "m { x = 1 }; x [0]", "m { x = 1 }; x([0])"},
		{"block parameter after the block", "m { |a| a }; a [0]", "m { |a| a }; a([0])"},
		{"local assigned later", "x [0]; x = 1", "x([0]); x = 1"},
		{"instance variable is no local", "@x = 1; x [0]", "@x = 1; x([0])"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := testParseAST(t, test.input).String(0)
			want := testParseAST(t, test.sameAs).String(0)

			if got != want {
				t.Fatalf("%q parsed as:\n%s\nwant the tree of %q:\n%s", test.input, got, test.sameAs, want)
			}
		})
	}
}

func TestIndexAndCommandArgumentHaveDifferentTrees(t *testing.T) {
	index := testParseAST(t, "x[0]").String(0)
	call := testParseAST(t, "x([0])").String(0)

	if index == call {
		t.Fatalf("both spellings print as %q, so tree comparisons prove nothing", index)
	}
}

func TestDeclareLocalVariablesMakesSpacedLeftBracketAnIndex(t *testing.T) {
	p := New(lexer.New(lexer.NewInput("test.rb", "x [0]")))
	p.DeclareLocalVariables([]string{"x"})

	got := p.ParseAST().String(0)
	want := testParseAST(t, "x[0]").String(0)

	if got != want {
		t.Fatalf("parsed as:\n%s\nwant:\n%s", got, want)
	}
}
