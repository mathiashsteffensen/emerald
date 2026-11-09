package lexer

import "unicode/utf8"

func isWhiteSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// identifierStartSize returns the byte length of the identifier-start
// character at pos in content, or 0 when no identifier can start there.
// Like Ruby, every valid non-ASCII UTF-8 character is an identifier character.
func identifierStartSize(content string, pos int) int {
	if pos >= len(content) {
		return 0
	}

	ch := content[pos]
	if ch < utf8.RuneSelf {
		if isLetter(ch) {
			return 1
		}
		return 0
	}

	r, size := utf8.DecodeRuneInString(content[pos:])
	if r == utf8.RuneError && size == 1 {
		return 0
	}

	return size
}

// identifierPartSize is identifierStartSize that also accepts digits.
func identifierPartSize(content string, pos int) int {
	if pos < len(content) && isDigit(content[pos]) {
		return 1
	}

	return identifierStartSize(content, pos)
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
