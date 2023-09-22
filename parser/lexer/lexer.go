package lexer

import (
	"strings"
)

type Lexer struct {
	inputChan    chan *Input
	currentInput *Input
	Line         int
	Column       int
	position     int // current position in input (points to current char)
	nextPosition int
	currentChar  byte // current char under examination
	outputChan   chan Token
	lastEmitted  Token

	templateNesting uint8
}

func New(input *Input) *Lexer {
	l := &Lexer{
		inputChan:  make(chan *Input, 10),
		outputChan: make(chan Token, 500),
		Line:       1,
		Column:     -1,
	}

	l.Feed(input)

	l.readChar()

	return l
}

func (l *Lexer) Snapshot(token Token) string {
	start := token.Pos - 80
	end := token.Pos + 8

	if start < 0 {
		start = 0
	}

	if end > len(l.currentInput.content) {
		end = len(l.currentInput.content)
	}

	var buf strings.Builder

	buf.WriteString(l.currentInput.content[start:end])
	buf.WriteString("\n")

	for i := 0; i < token.Column-2; i++ {
		buf.WriteString(" ")
	}
	buf.WriteString("^")

	return buf.String()
}

func (l *Lexer) Feed(input *Input) {
	if l.currentInput == nil {
		l.currentInput = input
		return
	}

	l.inputChan <- input
}

func (l *Lexer) inTemplate() bool {
	return l.templateNesting != 0
}

func (l *Lexer) Run() {
	go func() {
		var tok Token

		for tok.Type != EOF {

			l.eatWhitespace()

			switch l.currentChar {
			case '\n':
				l.Column = 1
				l.Line += 1
				tok = l.newToken(NEWLINE, l.currentChar)
			case '#':
				for l.currentChar != '\n' {
					l.readChar()
				}
				continue
			case '=':
				switch l.peekChar() {
				case '=':
					first := string(l.currentChar)
					l.readChar()

					if l.peekChar() == '=' {
						first += string(l.currentChar)
						l.readChar()
						tok = l.newTokenStr(CASE_EQ, first+string(l.currentChar))
					} else {
						tok = l.newTokenStr(EQ, first+string(l.currentChar))
					}
				case '>':
					tok = l.combineCurrentAndPeek(ARROW)
				case '~':
					tok = l.combineCurrentAndPeek(MATCH)
				default:
					tok = l.newToken(ASSIGN, l.currentChar)
				}
			case '!':
				if l.peekChar() == '=' {
					tok = l.combineCurrentAndPeek(NOT_EQ)
				} else {
					tok = l.newToken(BANG, l.currentChar)
				}
			case '?':
				tok = l.newToken(QUESTION, l.currentChar)
			case '<':
				if l.peekChar() == '=' {
					char := l.currentChar
					l.readChar()

					if l.peekChar() == '>' {
						tok = l.newTokenStr(SPACESHIP, string(char)+string(l.currentChar)+string(l.peekChar()))
						l.readChar()
					} else {
						tok = l.newTokenStr(LT_OR_EQ, string(char)+string(l.currentChar))
					}
				} else {
					if l.peekChar() == '<' {
						char := l.currentChar
						l.readChar()
						tok = l.newTokenStr(APPEND, string(char)+string(l.currentChar))
					} else {
						tok = l.newToken(LT, l.currentChar)
					}
				}
			case '>':
				if l.peekChar() == '=' {
					char := l.currentChar
					l.readChar()
					tok = Token{Type: GT_OR_EQ, Literal: string(char) + string(l.currentChar)}
				} else {
					tok = l.newToken(GT, l.currentChar)
				}
			case ';':
				tok = l.newToken(SEMICOLON, l.currentChar)
			case '(':
				tok = l.newToken(LPAREN, l.currentChar)
			case ')':
				tok = l.newToken(RPAREN, l.currentChar)
			case ',':
				tok = l.newToken(COMMA, l.currentChar)
			case '+':
				if l.peekChar() == '=' {
					tok = l.combineCurrentAndPeek(PLUS_ASSIGN)
				} else {
					tok = l.newToken(PLUS, l.currentChar)
				}
			case '-':
				if l.peekChar() == '=' {
					tok = l.combineCurrentAndPeek(MINUS_ASSIGN)
				} else {
					tok = l.newToken(MINUS, l.currentChar)
				}
			case '/':
				if l.isRegexpStart() {
					pattern := l.readRegexp()
					tok = l.newTokenStr(REGEXP, pattern)
				} else {
					if l.peekChar() == '=' {
						tok = l.combineCurrentAndPeek(SLASH_ASSIGN)
					} else {
						tok = l.newToken(SLASH, l.currentChar)
					}
				}
			case '*':
				if l.peekChar() == '=' {
					tok = l.combineCurrentAndPeek(ASTERISK_ASSIGN)
				} else {
					tok = l.newToken(ASTERISK, l.currentChar)
				}
			case '{':
				tok = l.newToken(LBRACE, l.currentChar)
			case '}':
				tok = l.newToken(RBRACE, l.currentChar)

				// If we are in a template i.e.
				// "This is a #{template <We are here>}"
				if l.inTemplate() {
					// Emit the RBRACE token immediately
					l.sendToken(tok)

					// Returns true if we are about to start another template
					// If we aren't we need to make sure to send ending string token
					if !l.lexDoubleQuotedString(&tok) {
						// Slight optimization, we only send the token if the ending string has characters
						//
						// This ending string doesn't have characters, so we end up only joining 2 strings
						// "This is a #{template <We are here>}"
						//
						// This ending string has characters, so we have to join 3 strings
						// "This is a #{template <We are here>} blah blah"
						if tok.Literal != "" {
							l.sendToken(tok)
						}

						l.readChar()
					}

					l.templateNesting -= 1

					continue
				}
			case '[':
				tok = l.newToken(LBRACKET, l.currentChar)
			case ']':
				tok = l.newToken(RBRACKET, l.currentChar)
			case ':':
				if l.peekChar() == ':' {
					// Scope operator
					// File::Stat
					char := l.currentChar
					l.readChar()
					tok = Token{Type: SCOPE, Literal: string(char) + string(l.currentChar)}
				} else if isLetter(l.peekChar()) {
					// A "normal" symbol
					// :symbol
					char := l.currentChar
					tok.Pos = l.position
					tok.Column = l.Column
					tok.Line = l.Line
					tok.Type = SYMBOL

					l.readChar()
					tok.Literal = string(char) + l.readIdentifier()
					l.sendToken(tok)

					continue
				} else if peek := l.peekChar(); peek == '"' || peek == '\'' {
					// A quoted symbol (does not support interpolation yet)
					// :"symbol" or :'symbol'
					char := l.currentChar
					tok.Pos = l.position
					tok.Column = l.Column
					tok.Line = l.Line
					tok.Type = SYMBOL

					l.readChar()
					tok.Literal = string(char) + l.readString(peek)
					l.readChar()
					l.sendToken(tok)

					continue
				} else {
					tok = l.newToken(COLON, l.currentChar)
				}
			case '.':
				if l.peekChar() == '.' {
					tok = l.combineCurrentAndPeek(RANGE_INCLUSIVE)

					if l.peekChar() == '.' {
						tok.Type = RANGE_EXCLUSIVE
						tok.Literal = tok.Literal + string(l.currentChar)
						l.readChar()
					}
				} else {
					tok = l.newToken(DOT, l.currentChar)
				}
			case '"':
				if l.lexDoubleQuotedString(&tok) {
					continue
				}
			case '\'':
				l.lexSingleQuotedString(&tok)
			case '&':
				if l.peekChar() == '&' {
					char := l.currentChar
					l.readChar()

					if l.peekChar() == '=' {
						secondChar := l.currentChar
						l.readChar()
						tok = Token{Type: BOOL_AND_ASSIGN, Literal: string(char) + string(secondChar) + string(l.currentChar)}
					} else {
						tok = Token{Type: BOOL_AND, Literal: string(char) + string(l.currentChar)}
					}
				} else {
					tok = l.newToken(BIT_AND, l.currentChar)
				}
			case '|':
				if l.peekChar() == '|' {
					char := l.currentChar
					l.readChar()

					if l.peekChar() == '=' {
						secondChar := l.currentChar
						l.readChar()
						tok = Token{Type: BOOL_OR_ASSIGN, Literal: string(char) + string(secondChar) + string(l.currentChar)}
					} else {
						tok = Token{Type: BOOL_OR, Literal: string(char) + string(l.currentChar)}
					}
				} else {
					tok = l.newToken(BIT_OR, l.currentChar)
				}
			case '@':
				tok.Type = INSTANCE_VAR

				char := l.currentChar

				l.readChar()

				tok.Literal = string(char) + l.readIdentifier()

				l.sendToken(tok)
				continue
			case '$':
				tok.Type = GLOBAL_IDENT

				char := l.currentChar

				l.readChar()

				var ident string

				if !isLetter(l.currentChar) {
					ident = string(char) + string(l.currentChar)
					l.readChar()
				} else {
					ident = string(char) + l.readIdentifier()
				}

				tok.Literal = ident

				l.sendToken(tok)
				continue
			case 0:
				tok.Literal = ""
				tok.Type = EOF
			default:
				if isLetter(l.currentChar) {
					tok.Pos = l.position
					tok.Column = l.Column
					tok.Line = l.Line
					tok.Literal = l.readIdentifier()
					tok.Type = lookupIdent(tok.Literal)
					l.sendToken(tok)

					l.eatWhitespace()

					if tok.Type == DEF && l.peekChar() == '@' && (l.currentChar == '-' || l.currentChar == '!') {
						l.sendToken(l.combineCurrentAndPeek(IDENT))
						l.readChar()
					}

					continue
				} else if isDigit(l.currentChar) {
					l.lexNumber(&tok)
					continue
				} else {
					tok = l.newToken(ILLEGAL, l.currentChar)
				}
			}

			l.readChar()

			l.sendToken(tok)
		}

		l.sendToken(tok)
	}()
}

func (l *Lexer) Close() {
	close(l.inputChan)
	close(l.outputChan)
}

func (l *Lexer) sendToken(token Token) {
	l.lastEmitted = token
	l.outputChan <- token
}

func (l *Lexer) NextToken() Token {
	return <-l.outputChan
}

func (l *Lexer) newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: l.Line, Column: l.Column, Pos: l.position}
}

func (l *Lexer) newTokenStr(tokenType TokenType, str string) Token {
	return Token{Type: tokenType, Literal: str, Line: l.Line, Column: l.Column, Pos: l.position}
}

func (l *Lexer) combineCurrentAndPeek(typ TokenType) Token {
	char := l.currentChar
	l.readChar()
	return l.newTokenStr(typ, string(char)+string(l.currentChar))
}

func (l *Lexer) readChar() {
	l.currentChar = l.peekChar()

	l.position = l.nextPosition
	l.nextPosition += 1

	l.Column += 1
}

func (l *Lexer) peekChar() byte {
	if l.currentInput == nil || l.nextPosition >= len(l.currentInput.content) {
		select {
		case nextInput := <-l.inputChan:
			l.currentInput = nextInput

			l.position = 0
			l.nextPosition = 1

			return 0
		default:
			return 0
		}
	} else {
		return l.currentInput.content[l.nextPosition]
	}
}

func (l *Lexer) eatWhitespace() {
	for isWhiteSpace(l.currentChar) {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.currentChar) || l.currentChar == '_' {
		l.readChar()
	}
	return l.currentInput.content[position:l.position]
}

func (l *Lexer) readRegexp() string {
	position := l.position + 1
	for {
		prevChar := l.currentChar
		l.readChar()
		if (l.currentChar == '/' && prevChar != '\\') || l.currentChar == 0 {
			break
		}
	}
	return l.currentInput.content[position:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.currentChar) || isDigit(l.currentChar) {
		l.readChar()
	}

	if l.currentChar == '?' || l.currentChar == '!' {
		l.readChar()
	}

	return l.currentInput.content[position:l.position]
}
