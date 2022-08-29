package lexer

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) isRegexpStart() bool {
	if l.lastEmitted.Type == "" {
		return true
	}

	switch l.lastEmitted.Type {
	case IDENT, INSTANCE_VAR, GLOBAL_IDENT, NULL, FALSE, TRUE, SELF, DEF, RBRACE, RBRACKET, INT, FLOAT, DOT:
		return false
	default:
		return true
	}
}
