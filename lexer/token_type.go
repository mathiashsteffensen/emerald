package lexer

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT        TokenType = "IDENT"
	INT          TokenType = "INT"
	STRING       TokenType = "STRING"
	INSTANCE_VAR TokenType = "INSTANCE_VAR"

	// Operators
	ASSIGN          TokenType = "="
	PLUS            TokenType = "+"
	MINUS           TokenType = "-"
	BANG            TokenType = "!"
	ASTERISK        TokenType = "*"
	SLASH           TokenType = "/"
	LT              TokenType = "<"
	GT              TokenType = ">"
	EQ              TokenType = "=="
	NOT_EQ          TokenType = "!="
	LT_OR_EQ        TokenType = "<="
	GT_OR_EQ        TokenType = ">="
	APPEND          TokenType = "<<"
	DOT             TokenType = "."
	BIT_AND         TokenType = "&"
	BOOL_AND        TokenType = "&&"
	BOOL_AND_ASSIGN TokenType = "&&="
	BIT_OR          TokenType = "|"
	BOOL_OR         TokenType = "||"
	BOOL_OR_ASSIGN  TokenType = "||="

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"

	// Keywords
	CLASS  TokenType = "CLASS"
	DEF    TokenType = "DEF"
	DO     TokenType = "DO"
	END    TokenType = "END"
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"
	SELF   TokenType = "SELF"
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	RETURN TokenType = "RETURN"
	NULL   TokenType = "NULL"
)

var keywords = map[string]TokenType{
	"class":  CLASS,
	"def":    DEF,
	"do":     DO,
	"true":   TRUE,
	"false":  FALSE,
	"self":   SELF,
	"if":     IF,
	"else":   ELSE,
	"end":    END,
	"return": RETURN,
	"nil":    NULL,
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
