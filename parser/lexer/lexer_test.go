package lexer

import (
	"strconv"
	"testing"
	"time"
)

func TestLexer_Run(t *testing.T) {
	input := `/r/
	five = 5.0;
	ten = 1_0;

	class Integer
		def -@(y)
			self + y
		end

		def !@ unless
	end

	@result = add!(five, :ten);

	@result.method?

	print(result)
	!-*5; case when
	5 < 10 > 5 / 5;

	[] << 5
	[].each { |w| do_stuff(w) }
	[].each do |w|
		do_stuff(w)
	end

	# This is a comment

	if (5 < 10)
		return true;
	elsif true
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
	nil &&= ||= module ?
	begin rescue ensure
	$: $LOAD_PATH Integer::Math::MAX <=> ===
	/^[w]|abc|\// =~ yield while break += -= ident /= *=
	#{this is a comment}
	"this is a #{template}" "this is a #{template} also" "this is a #{template} also #{boop.method} #{"nested #{template}"}" ½
	:"symbol" 'string' :'symbol' .. ...`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{REGEXP, "r"},
		{NEWLINE, "\n"},
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
		{IDENT, "-@"},
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
		{NEWLINE, "\n"},
		{DEF, "def"},
		{IDENT, "!@"},
		{UNLESS, "unless"},
		{NEWLINE, "\n"},
		{END, "end"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{INSTANCE_VAR, "@result"},
		{ASSIGN, "="},
		{IDENT, "add!"},
		{LPAREN, "("},
		{IDENT, "five"},
		{COMMA, ","},
		{SYMBOL, ":ten"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{INSTANCE_VAR, "@result"},
		{DOT, "."},
		{IDENT, "method?"},
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
		{CASE, "case"},
		{WHEN, "when"},
		{NEWLINE, "\n"},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{GT, ">"},
		{INT, "5"},
		{SLASH, "/"},
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
		{ELSIF, "elsif"},
		{TRUE, "true"},
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
		{QUESTION, "?"},
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
		{CASE_EQ, "==="},
		{NEWLINE, "\n"},
		{REGEXP, `^[w]|abc|\/`},
		{MATCH, "=~"},
		{YIELD, "yield"},
		{WHILE, "while"},
		{BREAK, "break"},
		{PLUS_ASSIGN, "+="},
		{MINUS_ASSIGN, "-="},
		{IDENT, "ident"},
		{SLASH_ASSIGN, "/="},
		{ASTERISK_ASSIGN, "*="},
		{NEWLINE, "\n"},
		{NEWLINE, "\n"},
		{STRING, "this is a "},
		{LTEMPLATE, "#{"},
		{IDENT, "template"},
		{RBRACE, "}"},
		{STRING, "this is a "},
		{LTEMPLATE, "#{"},
		{IDENT, "template"},
		{RBRACE, "}"},
		{STRING, " also"},
		{STRING, "this is a "},
		{LTEMPLATE, "#{"},
		{IDENT, "template"},
		{RBRACE, "}"},
		{STRING, " also "},
		{LTEMPLATE, "#{"},
		{IDENT, "boop"},
		{DOT, "."},
		{IDENT, "method"},
		{RBRACE, "}"},
		{STRING, " "},
		{LTEMPLATE, "#{"},
		{STRING, "nested "},
		{LTEMPLATE, "#{"},
		{IDENT, "template"},
		{RBRACE, "}"},
		{RBRACE, "}"},
		{ILLEGAL, "Â"},
		{ILLEGAL, "½"},
		{NEWLINE, "\n"},
		{SYMBOL, ":symbol"},
		{STRING, "string"},
		{SYMBOL, ":symbol"},
		{RANGE_INCLUSIVE, ".."},
		{RANGE_EXCLUSIVE, "..."},
		{EOF, ""},
	}

	l := New(NewInput("test.rb", input))

	l.Run()

	defer l.Close()

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			select {
			case tok := <-l.outputChan:
				testToken(t, tok, tt.expectedLiteral, tt.expectedType)
			case <-time.After(100 * time.Millisecond):
				t.Fatalf("Failed to receive token after 100ms, expected=%q (%q)", tt.expectedLiteral, tt.expectedType)
			}
		})
	}
}

func TestLexer_Feed(t *testing.T) {

}

func TestLexer_Snapshot(t *testing.T) {
	input := "w + w"

	l := New(NewInput("test.rb", input))

	l.Run()

	tok := l.NextToken()

	testToken(t, tok, "w", IDENT)

	l.NextToken()
	tok = l.NextToken()

	snapshot := l.Snapshot(tok)

	if snapshot == "w + w" {
		t.Fatalf("Snapshot was just source code without annotation")
	}
}

func testToken(t *testing.T, tok Token, literal string, typ TokenType) {
	if tok.Type != typ {
		t.Fatalf("tokentype wrong. expected=%q, got=%q (%q) ", typ, tok.Literal, tok.Type)
	}
	if tok.Literal != literal {
		t.Fatalf("literal wrong. expected=%q, got=%q", literal, tok.Literal)
	}
}
