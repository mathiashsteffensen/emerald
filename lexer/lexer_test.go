package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `five = 5;
	ten = 10;

	class Integer
		def add_num(y)
			self + y
		end
	end

	@result = add(five, ten);

	@result.method

	print(result)
	!-/*5;
	5 < 10 > 5;

	[] << 5
	[].each { |w| do_stuff(w) }
	[].each do |w|
		do_stuff(w)
	end

	# This is a comment

	if (5 < 10)
		return true;
	else
		return false;
	end

	10 == 10;
	10 != 9;
	10 >= 9
	9 <= 10
	"foobar"
	"foo bar"
	{hello: "world"}
	{hello: "world",}
	obj.hello.world
	obj["hello"][0][value]
	[0, 1]
	nil`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IDENT, "five"},
		{ASSIGN, "="},
		{INT, "5"},
		{SEMICOLON, ";"},
		{IDENT, "ten"},
		{ASSIGN, "="},
		{INT, "10"},
		{SEMICOLON, ";"},
		{CLASS, "class"},
		{IDENT, "Integer"},
		{DEF, "def"},
		{IDENT, "add_num"},
		{LPAREN, "("},
		{IDENT, "y"},
		{RPAREN, ")"},
		{SELF, "self"},
		{PLUS, "+"},
		{IDENT, "y"},
		{END, "end"},
		{END, "end"},
		{INSTANCE_VAR, "@result"},
		{ASSIGN, "="},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "five"},
		{COMMA, ","},
		{IDENT, "ten"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{INSTANCE_VAR, "@result"},
		{DOT, "."},
		{IDENT, "method"},
		{IDENT, "print"},
		{LPAREN, "("},
		{IDENT, "result"},
		{RPAREN, ")"},
		{BANG, "!"},
		{MINUS, "-"},
		{SLASH, "/"},
		{ASTERISK, "*"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{GT, ">"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{APPEND, "<<"},
		{INT, "5"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{DOT, "."},
		{IDENT, "each"},
		{LBRACE, "{"},
		{LINE, "|"},
		{IDENT, "w"},
		{LINE, "|"},
		{IDENT, "do_stuff"},
		{LPAREN, "("},
		{IDENT, "w"},
		{RPAREN, ")"},
		{RBRACE, "}"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{DOT, "."},
		{IDENT, "each"},
		{DO, "do"},
		{LINE, "|"},
		{IDENT, "w"},
		{LINE, "|"},
		{IDENT, "do_stuff"},
		{LPAREN, "("},
		{IDENT, "w"},
		{RPAREN, ")"},
		{END, "end"},
		{IF, "if"},
		{LPAREN, "("},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{RPAREN, ")"},
		{RETURN, "return"},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{ELSE, "else"},
		{RETURN, "return"},
		{FALSE, "false"},
		{SEMICOLON, ";"},
		{END, "end"},
		{INT, "10"},
		{EQ, "=="},
		{INT, "10"},
		{SEMICOLON, ";"},
		{INT, "10"},
		{NOT_EQ, "!="},
		{INT, "9"},
		{SEMICOLON, ";"},
		{INT, "10"},
		{GT_OR_EQ, ">="},
		{INT, "9"},
		{INT, "9"},
		{LT_OR_EQ, "<="},
		{INT, "10"},
		{STRING, "foobar"},
		{STRING, "foo bar"},
		{LBRACE, "{"},
		{IDENT, "hello"},
		{COLON, ":"},
		{STRING, "world"},
		{RBRACE, "}"},
		{LBRACE, "{"},
		{IDENT, "hello"},
		{COLON, ":"},
		{STRING, "world"},
		{COMMA, ","},
		{RBRACE, "}"},
		{IDENT, "obj"},
		{DOT, "."},
		{IDENT, "hello"},
		{DOT, "."},
		{IDENT, "world"},
		{IDENT, "obj"},
		{LBRACKET, "["},
		{STRING, "hello"},
		{RBRACKET, "]"},
		{LBRACKET, "["},
		{INT, "0"},
		{RBRACKET, "]"},
		{LBRACKET, "["},
		{IDENT, "value"},
		{RBRACKET, "]"},
		{LBRACKET, "["},
		{INT, "0"},
		{COMMA, ","},
		{INT, "1"},
		{RBRACKET, "]"},
		{NULL, "nil"},
		{EOF, ""},
	}

	l := New(NewInput("test.rb", input))

	l.Run()

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
