package lexer

import (
	"bytes"
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
}

func New(input *Input) *Lexer {
	l := &Lexer{inputChan: make(chan *Input, 100), outputChan: make(chan Token, 100), Line: 1, Column: -1}

	l.Feed(input)

	l.readChar()

	return l
}

func (l *Lexer) Snapshot(token Token) string {
	start := token.Pos - 30
	end := token.Pos + 5

	if start < 0 {
		start = 0
	}

	if end > len(l.currentInput.content) {
		end = len(l.currentInput.content)
	}

	var buf bytes.Buffer

	buf.WriteString(l.currentInput.content[start:end])
	buf.WriteString("\n")

	for i := 0; i < token.Pos-start; i++ {
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

func (l *Lexer) Run() {
	go func() {
		var tok Token

		for tok.Type != EOF {

			l.eatWhitespace()

			switch l.currentChar {
			case '#':
				for l.currentChar != '\n' {
					l.readChar()
				}
				continue
			case '=':
				if l.peekChar() == '=' {
					char := l.currentChar
					l.readChar()
					tok = Token{Type: EQ, Literal: string(char) + string(l.currentChar)}
				} else {
					tok = l.newToken(ASSIGN, l.currentChar)
				}
			case '!':
				if l.peekChar() == '=' {
					char := l.currentChar
					l.readChar()
					tok = Token{Type: NOT_EQ, Literal: string(char) + string(l.currentChar)}
				} else {
					tok = l.newToken(BANG, l.currentChar)
				}
			case '<':
				if l.peekChar() == '=' {
					char := l.currentChar
					l.readChar()
					tok = Token{Type: LT_OR_EQ, Literal: string(char) + string(l.currentChar)}
				} else {
					if l.peekChar() == '<' {
						char := l.currentChar
						l.readChar()
						tok = Token{Type: APPEND, Literal: string(char) + string(l.currentChar)}
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
				tok = l.newToken(PLUS, l.currentChar)
			case '-':
				tok = l.newToken(MINUS, l.currentChar)
			case '/':
				tok = l.newToken(SLASH, l.currentChar)
			case '*':
				tok = l.newToken(ASTERISK, l.currentChar)
			case '{':
				tok = l.newToken(LBRACE, l.currentChar)
			case '}':
				tok = l.newToken(RBRACE, l.currentChar)
			case '[':
				tok = l.newToken(LBRACKET, l.currentChar)
			case ']':
				tok = l.newToken(RBRACKET, l.currentChar)
			case ':':
				tok = l.newToken(COLON, l.currentChar)
			case '.':
				tok = l.newToken(DOT, l.currentChar)
			case '"':
				tok.Type = STRING
				tok.Literal = l.readString()
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
					continue
				} else if isDigit(l.currentChar) {
					tok.Pos = l.position
					tok.Column = l.Column
					tok.Line = l.Line
					tok.Type = INT
					tok.Literal = l.readNumber()
					l.sendToken(tok)
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
	l.outputChan <- token
}

func (l *Lexer) NextToken() Token {
	return <-l.outputChan
}

func (l *Lexer) newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: l.Line, Column: l.Column, Pos: l.position}
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
		default:
			return 0
		}
		return 0
	} else {
		return l.currentInput.content[l.nextPosition]
	}
}

func (l *Lexer) eatWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		isNewLine := l.currentChar == '\n'

		l.readChar()

		if isNewLine {
			l.Column = 1
			l.Line += 1
		}
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.currentChar) || l.currentChar == '_' {
		l.readChar()
	}
	return l.currentInput.content[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.currentChar == '"' || l.currentChar == 0 {
			break
		}
	}
	return l.currentInput.content[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.currentChar) {
		l.readChar()
	}
	return l.currentInput.content[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
