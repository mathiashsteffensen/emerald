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
	SYMBOL TokenType = "SYMBOL"
	REGEXP TokenType = "REGEXP"

	// Operators
	PLUS            TokenType = "+"
	MINUS           TokenType = "-"
	BANG            TokenType = "!"
	ASTERISK        TokenType = "*"
	SLASH           TokenType = "/"
	SPACESHIP       TokenType = "<=>"
	LT              TokenType = "<"
	GT              TokenType = ">"
	EQ              TokenType = "=="
	CASE_EQ         TokenType = "==="
	NOT_EQ          TokenType = "!="
	MATCH           TokenType = "=~"
	LT_OR_EQ        TokenType = "<="
	GT_OR_EQ        TokenType = ">="
	ASSIGN          TokenType = "="
	PLUS_ASSIGN     TokenType = "+="
	MINUS_ASSIGN    TokenType = "-="
	SLASH_ASSIGN    TokenType = "/="
	ASTERISK_ASSIGN TokenType = "*="
	BOOL_AND_ASSIGN TokenType = "&&="
	BOOL_OR_ASSIGN  TokenType = "||="
	APPEND          TokenType = "<<"
	DOT             TokenType = "."
	RANGE_INCLUSIVE TokenType = ".."
	RANGE_EXCLUSIVE TokenType = "..."
	SCOPE           TokenType = "::"
	BIT_AND         TokenType = "&"
	BOOL_AND        TokenType = "&&"
	BIT_OR          TokenType = "|"
	BOOL_OR         TokenType = "||"
	QUESTION        TokenType = "?"

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"
	ARROW     TokenType = "=>"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LTEMPLATE TokenType = "#{"
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
	YIELD  TokenType = "YIELD"
	IF     TokenType = "IF"
	ELSIF  TokenType = "ELSIF"
	UNLESS TokenType = "UNLESS"
	CASE   TokenType = "CASE"
	WHEN   TokenType = "WHEN"
	ELSE   TokenType = "ELSE"
	RETURN TokenType = "RETURN"
	NULL   TokenType = "NULL"
	WHILE  TokenType = "WHILE"
	BREAK  TokenType = "BREAK"
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
	"yield":  YIELD,
	"if":     IF,
	"elsif":  ELSIF,
	"unless": UNLESS,
	"case":   CASE,
	"when":   WHEN,
	"else":   ELSE,
	"end":    END,
	"return": RETURN,
	"nil":    NULL,
	"while":  WHILE,
	"break":  BREAK,
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
