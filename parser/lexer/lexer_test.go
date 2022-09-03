package lexer

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `five = 5.0;
	ten = 1_0;

	class Integer
		def add_num(y)
			self + y
		end
	end

	@result = add(five, ten);

	@result.method

	print(result)
	!-*5;
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
	9 <= 10 & 4 && 4 | 4 || 4
	"foobar"
	"foo bar"
	{hello: "world"}
	{hello => "world",}
	obj.hello.world
	obj["hello"][0][value]
	[0, 1]
	nil &&= ||= module
	begin rescue ensure
	$: $LOAD_PATH Integer::Math::MAX <=>
	/^[w]|abc/ =~ yield`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{IDENT, "five"},
		{ASSIGN, "="},
		{FLOAT, "5.0"},
		{SEMICOLON, ";"},
		{NEWLINE, "\n"},
		{IDENT, "ten"},
		{ASSIGN, "="},
		{INT, "1_0"},
		{SEMICOLON, ";"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{CLASS, "class"},
		{IDENT, "Integer"},
		{NEWLINE, "\n"},
		{DEF, "def"},
		{IDENT, "add_num"},
		{LPAREN, "("},
		{IDENT, "y"},
		{RPAREN, ")"},
		{NEWLINE, "\n"},
		{SELF, "self"},
		{PLUS, "+"},
		{IDENT, "y"},
		{NEWLINE, "\n"},
		{END, "end"},
		{NEWLINE, "\n"},
		{END, "end"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{INSTANCE_VAR, "@result"},
		{ASSIGN, "="},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "five"},
		{COMMA, ","},
		{IDENT, "ten"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{INSTANCE_VAR, "@result"},
		{DOT, "."},
		{IDENT, "method"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{IDENT, "print"},
		{LPAREN, "("},
		{IDENT, "result"},
		{RPAREN, ")"},
		{NEWLINE, "\n"},
		{BANG, "!"},
		{MINUS, "-"},
		{ASTERISK, "*"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{NEWLINE, "\n"},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{GT, ">"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{APPEND, "<<"},
		{INT, "5"},
		{NEWLINE, "\n"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{DOT, "."},
		{IDENT, "each"},
		{LBRACE, "{"},
		{BIT_OR, "|"},
		{IDENT, "w"},
		{BIT_OR, "|"},
		{IDENT, "do_stuff"},
		{LPAREN, "("},
		{IDENT, "w"},
		{RPAREN, ")"},
		{RBRACE, "}"},
		{NEWLINE, "\n"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{DOT, "."},
		{IDENT, "each"},
		{DO, "do"},
		{BIT_OR, "|"},
		{IDENT, "w"},
		{BIT_OR, "|"},
		{NEWLINE, "\n"},
		{IDENT, "do_stuff"},
		{LPAREN, "("},
		{IDENT, "w"},
		{RPAREN, ")"},
		{NEWLINE, "\n"},
		{END, "end"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{IF, "if"},
		{LPAREN, "("},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{RPAREN, ")"},
		{NEWLINE, "\n"},
		{RETURN, "return"},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{NEWLINE, "\n"},
		{ELSE, "else"},
		{NEWLINE, "\n"},
		{RETURN, "return"},
		{FALSE, "false"},
		{SEMICOLON, ";"},
		{NEWLINE, "\n"},
		{END, "end"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{INT, "10"},
		{EQ, "=="},
		{INT, "10"},
		{SEMICOLON, ";"},
		{NEWLINE, "\n"},
		{INT, "10"},
		{NOT_EQ, "!="},
		{INT, "9"},
		{SEMICOLON, ";"},
		{NEWLINE, "\n"},
		{INT, "10"},
		{GT_OR_EQ, ">="},
		{INT, "9"},
		{NEWLINE, "\n"},
		{INT, "9"},
		{LT_OR_EQ, "<="},
		{INT, "10"},
		{BIT_AND, "&"},
		{INT, "4"},
		{BOOL_AND, "&&"},
		{INT, "4"},
		{BIT_OR, "|"},
		{INT, "4"},
		{BOOL_OR, "||"},
		{INT, "4"},
		{NEWLINE, "\n"},
		{STRING, "foobar"},
		{NEWLINE, "\n"},
		{STRING, "foo bar"},
		{NEWLINE, "\n"},
		{LBRACE, "{"},
		{IDENT, "hello"},
		{COLON, ":"},
		{STRING, "world"},
		{RBRACE, "}"},
		{NEWLINE, "\n"},
		{LBRACE, "{"},
		{IDENT, "hello"},
		{ARROW, "=>"},
		{STRING, "world"},
		{COMMA, ","},
		{RBRACE, "}"},
		{NEWLINE, "\n"},
		{IDENT, "obj"},
		{DOT, "."},
		{IDENT, "hello"},
		{DOT, "."},
		{IDENT, "world"},
		{NEWLINE, "\n"},
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
		{NEWLINE, "\n"},
		{LBRACKET, "["},
		{INT, "0"},
		{COMMA, ","},
		{INT, "1"},
		{RBRACKET, "]"},
		{NEWLINE, "\n"},
		{NULL, "nil"},
		{BOOL_AND_ASSIGN, "&&="},
		{BOOL_OR_ASSIGN, "||="},
		{MODULE, "module"},
		{NEWLINE, "\n"},
		{BEGIN, "begin"},
		{RESCUE, "rescue"},
		{ENSURE, "ensure"},
		{NEWLINE, "\n"},
		{GLOBAL_IDENT, "$:"},
		{GLOBAL_IDENT, "$LOAD_PATH"},
		{IDENT, "Integer"},
		{SCOPE, "::"},
		{IDENT, "Math"},
		{SCOPE, "::"},
		{IDENT, "MAX"},
		{SPACESHIP, "<=>"},
		{NEWLINE, "\n"},
		{REGEXP, "^[w]|abc"},
		{MATCH, "=~"},
		{YIELD, "yield"},
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
