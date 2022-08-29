package lexer

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers
	IDENT        TokenType = "IDENT"
	GLOBAL_IDENT TokenType = "GLOBAL_IDENT"
	INSTANCE_VAR TokenType = "INSTANCE_VAR"

	// Literals
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	STRING TokenType = "STRING"
	REGEXP TokenType = "REGEXP"

	// Operators
	ASSIGN          TokenType = "="
	PLUS            TokenType = "+"
	MINUS           TokenType = "-"
	BANG            TokenType = "!"
	ASTERISK        TokenType = "*"
	SLASH           TokenType = "/"
	SPACESHIP       TokenType = "<=>"
	LT              TokenType = "<"
	GT              TokenType = ">"
	EQ              TokenType = "=="
	NOT_EQ          TokenType = "!="
	MATCH           TokenType = "=~"
	LT_OR_EQ        TokenType = "<="
	GT_OR_EQ        TokenType = ">="
	APPEND          TokenType = "<<"
	DOT             TokenType = "."
	SCOPE           TokenType = "::"
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
	ARROW     TokenType = "=>"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LBRACKET  TokenType = "["
	RBRACKET  TokenType = "]"
	NEWLINE   TokenType = "\n"

	// Keywords
	CLASS  TokenType = "CLASS"
	MODULE TokenType = "MODULE"
	DEF    TokenType = "DEF"
	DO     TokenType = "DO"
	BEGIN  TokenType = "BEGIN"
	RESCUE TokenType = "RESCUE"
	ENSURE TokenType = "ENSURE"
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
	"module": MODULE,
	"def":    DEF,
	"do":     DO,
	"begin":  BEGIN,
	"rescue": RESCUE,
	"ensure": ENSURE,
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
