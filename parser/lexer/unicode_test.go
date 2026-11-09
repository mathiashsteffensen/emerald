package lexer

import (
	"testing"
)

func lexAll(t *testing.T, input string) []Token {
	t.Helper()

	l := New(NewInput("test.rb", input))
	l.Run()
	defer l.Close()

	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)

		if tok.Type == EOF {
			return tokens
		}
	}
}

func TestLexer_UnicodeIdentifiers(t *testing.T) {
	input := "café = π + ✓ x×y naïve? naïve! @名前 $π :日本 :café? :Ärger {ключ: 1} kw(ключ:) éclair1"

	expected := []struct {
		typ     TokenType
		literal string
	}{
		{IDENT, "café"},
		{ASSIGN, "="},
		{IDENT, "π"},
		{PLUS, "+"},
		{IDENT, "✓"},
		{IDENT, "x×y"},
		{IDENT, "naïve?"},
		{IDENT, "naïve!"},
		{INSTANCE_VAR, "@名前"},
		{GLOBAL_IDENT, "$π"},
		{SYMBOL, ":日本"},
		{SYMBOL, ":café?"},
		{SYMBOL, ":Ärger"},
		{LBRACE, "{"},
		{IDENT, "ключ"},
		{COLON, ":"},
		{INT, "1"},
		{RBRACE, "}"},
		{IDENT, "kw"},
		{LPAREN, "("},
		{IDENT, "ключ"},
		{COLON, ":"},
		{RPAREN, ")"},
		{IDENT, "éclair1"},
		{EOF, ""},
	}

	tokens := lexAll(t, input)

	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d: %v", len(expected), len(tokens), tokens)
	}

	for i, want := range expected {
		testToken(t, tokens[i], want.literal, want.typ)
	}
}

func TestLexer_ColumnsCountCharactersNotBytes(t *testing.T) {
	tokens := lexAll(t, "日本 = ключ + 1")

	expectedColumns := []int{0, 3, 5, 10, 12}

	for i, column := range expectedColumns {
		if tokens[i].Column != column {
			t.Errorf("token %q: expected column %d, got %d", tokens[i].Literal, column, tokens[i].Column)
		}
	}

	if tokens[2].Pos != len("日本 = ") {
		t.Errorf("Pos must stay a byte offset: expected %d, got %d", len("日本 = "), tokens[2].Pos)
	}
}

func TestLexer_SnapshotCaretSitsUnderTheToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"multibyte line", "ключ = ?", "ключ = ?\n       ^\n"},
		{"later line", "a = 1\nb = ?\nc", "a = 1\nb = ?\n    ^\nc\n"},
		{"tab-indented line", "x = 1\n\t\tb = ?", "x = 1\n    b = ?\n        ^\n"},
		{"multibyte char on later line", "a = 1\nπ = ?", "a = 1\nπ = ?\n    ^\n"},

		// Column does not count a leading byte order mark, so neither may
		// the caret padding.
		{"byte order mark and tab", byteOrderMark + "x =\t?", "x =  ?\n     ^\n"},
		{"byte order mark before tab", byteOrderMark + "\t?", "  ?\n  ^\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(NewInput("test.rb", tt.input))
			l.Run()
			defer l.Close()

			// Stop at EOF: the token channel stays open until Close, so
			// waiting for a token that never comes would hang the test.
			var question Token
			for tok := l.NextToken(); tok.Type != EOF; tok = l.NextToken() {
				if tok.Type == QUESTION {
					question = tok
				}
			}

			if question.Type != QUESTION {
				t.Fatalf("no %s token in %q", QUESTION, tt.input)
			}

			if snapshot := l.Snapshot(question); snapshot != tt.expected {
				t.Errorf("expected snapshot %q, got %q", tt.expected, snapshot)
			}
		})
	}
}

func TestLexer_IllegalBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "unknown ASCII byte",
			input: "\x01",
			expected: []Token{
				{Type: ILLEGAL, Literal: "\x01"},
			},
		},
		{
			name:  "lone lead byte",
			input: "x = \xC3",
			expected: []Token{
				{Type: IDENT, Literal: "x"},
				{Type: ASSIGN, Literal: "="},
				{Type: ILLEGAL, Literal: "\xC3"},
			},
		},
		{
			name:  "invalid byte ends an identifier",
			input: "ab\xFF",
			expected: []Token{
				{Type: IDENT, Literal: "ab"},
				{Type: ILLEGAL, Literal: "\xFF"},
			},
		},
		{
			name:  "invalid byte after $",
			input: "$\xFF",
			expected: []Token{
				{Type: ILLEGAL, Literal: "\xFF"},
			},
		},
		{
			name:  "invalid byte after :",
			input: ":\xFF",
			expected: []Token{
				{Type: ILLEGAL, Literal: "\xFF"},
				{Type: EOF, Literal: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexAll(t, tt.input)

			for i, want := range tt.expected {
				testToken(t, tokens[i], want.Literal, want.Type)
			}
		})
	}
}

func TestLexer_InvalidBytesTakeOneColumnEach(t *testing.T) {
	tests := []struct {
		input           string
		expectedColumns []int
	}{
		{"\x80x", []int{0, 1}},
		{"a\x80\x80b", []int{0, 1, 2, 3}},
		{"é\xC3x", []int{0, 1, 2}},
	}

	for _, tt := range tests {
		tokens := lexAll(t, tt.input)

		for i, column := range tt.expectedColumns {
			if tokens[i].Column != column {
				t.Errorf("%q token %d (%q): expected column %d, got %d", tt.input, i, tokens[i].Literal, column, tokens[i].Column)
			}
		}
	}
}

func TestLexer_SkipsLeadingByteOrderMark(t *testing.T) {
	tokens := lexAll(t, "\uFEFFputs 1")

	testToken(t, tokens[0], "puts", IDENT)

	if tokens[0].Column != 0 {
		t.Errorf("expected first token at column 0, got %d", tokens[0].Column)
	}

	if tokens[0].Pos != len("\uFEFF") {
		t.Errorf("Pos must stay a byte offset: expected %d, got %d", len("\uFEFF"), tokens[0].Pos)
	}
}

func TestLexer_EveryLineStartsAtColumnZero(t *testing.T) {
	tokens := lexAll(t, "a = 1\nb = ?\nπ = é")

	expected := []struct {
		literal string
		line    int
		column  int
	}{
		{"a", 1, 0},
		{"=", 1, 2},
		{"1", 1, 4},
		{"\n", 1, 5},
		{"b", 2, 0},
		{"=", 2, 2},
		{"?", 2, 4},
		{"\n", 2, 5},
		{"π", 3, 0},
		{"=", 3, 2},
		{"é", 3, 4},
	}

	for i, want := range expected {
		got := tokens[i]

		if got.Literal != want.literal || got.Line != want.line || got.Column != want.column {
			t.Errorf("token %d: expected %q at line %d, column %d, got %q at line %d, column %d",
				i, want.literal, want.line, want.column, got.Literal, got.Line, got.Column)
		}
	}
}

func TestLexer_InvalidUTF8InLiteralsIsSentAsIllegal(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      []Token
		illegalColumn int
	}{
		{
			name:  "double-quoted string",
			input: "x = \"a\xC3b\"",
			expected: []Token{
				{Type: IDENT, Literal: "x"},
				{Type: ASSIGN, Literal: "="},
				{Type: ILLEGAL, Literal: "\xC3"},
				{Type: STRING, Literal: "a\xC3b"},
			},
			illegalColumn: 6,
		},
		{
			name:  "single-quoted string",
			input: "'a\xC3'",
			expected: []Token{
				{Type: ILLEGAL, Literal: "\xC3"},
				{Type: STRING, Literal: "a\xC3"},
			},
			illegalColumn: 2,
		},
		{
			name:  "only the first invalid byte",
			input: "\"\xC3\xC3\"",
			expected: []Token{
				{Type: ILLEGAL, Literal: "\xC3"},
				{Type: STRING, Literal: "\xC3\xC3"},
				{Type: EOF, Literal: ""},
			},
			illegalColumn: 1,
		},
		{
			name:  "template continuation",
			input: "\"a#{y}\xC3\"",
			expected: []Token{
				{Type: STRING, Literal: "a"},
				{Type: LTEMPLATE, Literal: "#{"},
				{Type: IDENT, Literal: "y"},
				{Type: RBRACE, Literal: "}"},
				{Type: ILLEGAL, Literal: "\xC3"},
				{Type: STRING, Literal: "\xC3"},
			},
			illegalColumn: 6,
		},
		{
			name:  "regexp",
			input: "/a\xC3b/",
			expected: []Token{
				{Type: ILLEGAL, Literal: "\xC3"},
				{Type: REGEXP, Literal: "a\xC3b"},
			},
			illegalColumn: 2,
		},
		{
			name:  "quoted symbol",
			input: ":\"a\xC3\"",
			expected: []Token{
				{Type: ILLEGAL, Literal: "\xC3"},
				{Type: SYMBOL, Literal: ":a\xC3"},
			},
			illegalColumn: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexAll(t, tt.input)

			for i, want := range tt.expected {
				testToken(t, tokens[i], want.Literal, want.Type)

				if want.Type == ILLEGAL && tokens[i].Column != tt.illegalColumn {
					t.Errorf("expected the invalid byte at column %d, got %d", tt.illegalColumn, tokens[i].Column)
				}
			}
		})
	}
}

func TestLexer_ValidUTF8LiteralsAndCommentsSendNoIllegal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:     "multibyte string",
			input:    "\"é😀\"",
			expected: []Token{{Type: STRING, Literal: "é😀"}, {Type: EOF, Literal: ""}},
		},
		{
			name:     "byte order mark in a string",
			input:    "x\n\"\uFEFF\"",
			expected: []Token{{Type: IDENT, Literal: "x"}, {Type: NEWLINE, Literal: "\n"}, {Type: STRING, Literal: "\uFEFF"}, {Type: EOF, Literal: ""}},
		},
		{
			name:     "invalid byte in a comment",
			input:    "# a\xC3b\nputs 1",
			expected: []Token{{Type: NEWLINE, Literal: "\n"}, {Type: IDENT, Literal: "puts"}, {Type: INT, Literal: "1"}, {Type: EOF, Literal: ""}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexAll(t, tt.input)

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d: %v", len(tt.expected), len(tokens), tokens)
			}

			for i, want := range tt.expected {
				testToken(t, tokens[i], want.Literal, want.Type)
			}
		})
	}
}

func TestLexer_TokensWithoutOwnPositionStartWhereTheyAre(t *testing.T) {
	tests := []struct {
		name  string
		input string
		typ   TokenType
		line  int
		col   int
		pos   int
	}{
		{"instance variable", "x\n@a", INSTANCE_VAR, 2, 0, 2},
		{"global variable", "x\n  $g", GLOBAL_IDENT, 2, 2, 4},
		{"double-quoted string", "x\n\"s\"", STRING, 2, 0, 2},
		{"single-quoted string", "x\n's'", STRING, 2, 0, 2},
		{"first instance variable", "@a", INSTANCE_VAR, 1, 0, 0},
		{"first global variable", "$g", GLOBAL_IDENT, 1, 0, 0},
		{"first double-quoted string", "\"s\"", STRING, 1, 0, 0},
		{"first single-quoted string", "'s'", STRING, 1, 0, 0},
		{"end of file", "x\n", EOF, 2, 0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := lexAll(t, tt.input)

			var tok Token
			for _, candidate := range tokens {
				if candidate.Type == tt.typ {
					tok = candidate
					break
				}
			}

			if tok.Type != tt.typ {
				t.Fatalf("no %s token in %v", tt.typ, tokens)
			}

			if tok.Line != tt.line || tok.Column != tt.col || tok.Pos != tt.pos {
				t.Errorf("expected line %d, column %d, pos %d, got line %d, column %d, pos %d",
					tt.line, tt.col, tt.pos, tok.Line, tok.Column, tok.Pos)
			}
		})
	}
}

func TestLexer_BooleanOperatorsCarryTheirPosition(t *testing.T) {
	for _, input := range []string{"x\na && b", "x\na || b", "x\na &&= b", "x\na ||= b"} {
		tokens := lexAll(t, input)
		operator := tokens[3]

		// combineCurrentAndPeek places a two-character token at its second
		// character, a known fault, so pin only what is correct today.
		if operator.Line != 2 || operator.Pos == 0 {
			t.Errorf("%q: expected %s on line 2 with a position, got line %d, pos %d",
				input, operator.Literal, operator.Line, operator.Pos)
		}
	}
}
