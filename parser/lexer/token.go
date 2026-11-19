package lexer

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
	Pos     int

	// SpaceBefore is true when whitespace comes just before the token, as
	// in `puts [1]`. It tells an array argument from an index.
	SpaceBefore bool
}
