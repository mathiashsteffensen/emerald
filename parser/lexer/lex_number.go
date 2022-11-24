package lexer

func (l *Lexer) lexNumber(token *Token) {
	token.Pos = l.position
	token.Column = l.Column
	token.Line = l.Line
	token.Literal = l.readNumber()

	if l.currentChar == '.' && isDigit(l.peekChar()) {
		l.readChar()
		token.Literal += "." + l.readNumber()
		token.Type = FLOAT
	} else {
		token.Type = INT
	}

	if l.currentChar == 'e' || l.currentChar == 'E' {
		l.readChar()
		token.Literal += "e"
		if l.currentChar == '-' || l.currentChar == '+' {
			token.Literal += string(l.currentChar)
			l.readChar()
		}
		token.Literal += l.readNumber()
		token.Type = FLOAT
	}

	l.sendToken(*token)
}
