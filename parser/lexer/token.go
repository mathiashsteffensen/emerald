package lexer

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
	Pos     int
}
