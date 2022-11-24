package lexer

func (l *Lexer) lexDoubleQuotedString(tok *Token) bool {
	tok.Type = STRING
	tok.Literal = l.readString()

	if l.nextIsLTEMPLATE() {
		l.templateNesting += 1
		l.sendToken(*tok)
		l.sendToken(l.combineCurrentAndPeek(LTEMPLATE))
		l.readChar()
		return true
	}

	return false
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()

		if l.nextIsLTEMPLATE() {
			break
		}

		if l.currentChar == '"' || l.currentChar == 0 {
			break
		}
	}

	return l.currentInput.content[position:l.position]
}

func (l *Lexer) nextIsLTEMPLATE() bool {
	return l.currentChar == '#' && l.peekChar() == '{'
}
