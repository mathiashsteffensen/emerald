package lexer

import (
	"context"
	"strings"
	"sync"
	"unicode/utf8"
)

const byteOrderMark = "\uFEFF"

type Lexer struct {
	currentInput *Input
	Line         int
	Column       int
	position     int // current position in input (points to current char)
	nextPosition int

	// trailingBytesInChar counts the bytes of the current multibyte
	// character that readChar has yet to read. They add no column.
	trailingBytesInChar int

	// currentCharIsInvalidUTF8 is true when currentChar is a byte that does
	// not start a valid UTF-8 character.
	currentCharIsInvalidUTF8 bool

	currentChar byte // current char under examination
	outputChan  chan Token
	lastEmitted Token

	templateNesting uint8

	ctx       context.Context
	cancel    context.CancelFunc
	done      chan struct{}
	runOnce   sync.Once
	closeOnce sync.Once
}

func New(input *Input) *Lexer {
	return NewContext(context.Background(), input)
}

func NewContext(ctx context.Context, input *Input) *Lexer {
	if ctx == nil {
		ctx = context.Background()
	}

	l := &Lexer{
		outputChan: make(chan Token, 500),
		Line:       1,
		Column:     -1,
		ctx:        ctx,
		done:       make(chan struct{}),
	}

	l.currentInput = input

	// Like Ruby, skip a leading byte order mark.
	if strings.HasPrefix(input.content, byteOrderMark) {
		l.nextPosition = len(byteOrderMark)
	}

	l.readChar()

	return l
}

func (l *Lexer) Snapshot(token Token) string {
	// Column does not count a leading byte order mark, so drop it here too.
	lines := strings.Split(strings.TrimPrefix(l.currentInput.content, byteOrderMark), "\n")

	// token.Line is 1-based, convert to 0-based index
	tokenLineIdx := token.Line - 1

	// Determine the range of lines to include
	startLine := tokenLineIdx - 1
	if startLine < 0 {
		startLine = 0
	}

	endLine := tokenLineIdx + 1
	if endLine >= len(lines) {
		endLine = len(lines) - 1
	}

	var buf strings.Builder

	writeLine := func(lineNum int) {
		line := strings.ReplaceAll(lines[lineNum], "\t", "  ")
		buf.WriteString(line)
		buf.WriteString("\n")
	}

	// Add first line
	if startLine != tokenLineIdx {
		writeLine(startLine)
	}

	// Add second line (the one with the token)
	writeLine(tokenLineIdx)

	// Pad the caret the way writeLine shows the characters before the token:
	// a tab as two spaces, any other character, multibyte too, as one.
	for i, char := range []rune(lines[tokenLineIdx]) {
		if i >= token.Column {
			break
		}

		if char == '\t' {
			buf.WriteString("  ")
		} else {
			buf.WriteString(" ")
		}
	}
	buf.WriteString("^")

	buf.WriteString("\n")

	// Add third line
	if endLine != tokenLineIdx {
		writeLine(endLine)
	}

	return buf.String()
}

func (l *Lexer) inTemplate() bool {
	return l.templateNesting != 0
}

func (l *Lexer) Run() {
	l.RunContext(l.ctx)
}

func (l *Lexer) RunContext(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}

	l.runOnce.Do(func() {
		l.ctx, l.cancel = context.WithCancel(ctx)

		go l.run()
	})
}

func (l *Lexer) run() {
	defer close(l.done)

	var tok Token

	for tok.Type != EOF && !l.canceled() {
		l.eatWhitespace()

		if l.canceled() {
			return
		}

		// Tokens that set only their type and literal start here.
		tok = Token{Line: l.Line, Column: l.Column, Pos: l.position}

		switch l.currentChar {
		case '\n':
			// The NEWLINE token keeps its own place at the end of its line.
			// The next line starts at column 0, as line 1 does: readChar
			// moves the column from -1 to 0.
			tok = l.newToken(NEWLINE, l.currentChar)
			l.Line += 1
			l.Column = -1
		case '#':
			for l.currentChar != '\n' && l.currentChar != 0 && !l.canceled() {
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
				tok = l.combineCurrentAndPeek(GT_OR_EQ)
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
				tok = l.combineCurrentAndPeek(SCOPE)
			} else if l.identifierStartsAt(l.nextPosition) {
				// A "normal" symbol
				// :symbol
				char := l.currentChar
				tok.Type = SYMBOL

				l.readChar()
				tok.Literal = string(char) + l.readIdentifier()
				l.sendToken(tok)

				continue
			} else if peek := l.peekChar(); peek == '"' || peek == '\'' {
				// A quoted symbol (does not support interpolation yet)
				// :"symbol" or :'symbol'
				char := l.currentChar
				tok.Type = SYMBOL

				l.readChar()
				tok.Literal = string(char) + l.readString(peek)
				l.readChar()
				l.sendToken(tok)

				continue
			} else if l.invalidUTF8At(l.nextPosition) {
				// Invalid UTF-8 after `:` is no symbol name: send the byte
				// so that the parser reports it before anything else.
				l.readChar()
				tok = l.newToken(ILLEGAL, l.currentChar)
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
				tok = l.combineCurrentAndPeek(BOOL_AND)

				if l.peekChar() == '=' {
					l.readChar()
					tok.Type = BOOL_AND_ASSIGN
					tok.Literal += string(l.currentChar)
				}
			} else {
				tok = l.newToken(BIT_AND, l.currentChar)
			}
		case '|':
			if l.peekChar() == '|' {
				tok = l.combineCurrentAndPeek(BOOL_OR)

				if l.peekChar() == '=' {
					l.readChar()
					tok.Type = BOOL_OR_ASSIGN
					tok.Literal += string(l.currentChar)
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

			if l.identifierStartsAt(l.position) {
				ident = string(char) + l.readIdentifier()
			} else if l.currentChar >= utf8.RuneSelf {
				// Invalid UTF-8 after `$`: neither a name nor a special global like `$!`
				tok = l.newToken(ILLEGAL, l.currentChar)
				l.readChar()
				l.sendToken(tok)
				continue
			} else {
				ident = string(char) + string(l.currentChar)
				l.readChar()
			}

			tok.Literal = ident

			l.sendToken(tok)
			continue
		case 0:
			tok.Literal = ""
			tok.Type = EOF
		default:
			if l.identifierStartsAt(l.position) {
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

	if !l.canceled() {
		l.sendToken(tok)
	}
}

func (l *Lexer) Close() {
	l.closeOnce.Do(func() {
		notStarted := false
		l.runOnce.Do(func() {
			notStarted = true
			close(l.done)
		})

		if !notStarted {
			l.cancel()
			<-l.done
		}

		close(l.outputChan)
	})
}

func (l *Lexer) sendToken(token Token) {
	l.lastEmitted = token

	select {
	case l.outputChan <- token:
	case <-l.ctx.Done():
	}
}

func (l *Lexer) NextToken() Token {
	token, _ := l.NextTokenContext(l.ctx)
	return token
}

func (l *Lexer) NextTokenContext(ctx context.Context) (Token, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case token, ok := <-l.outputChan:
		if !ok {
			return Token{Type: EOF}, nil
		}
		return token, nil
	case <-ctx.Done():
		return Token{}, ctx.Err()
	}
}

func (l *Lexer) Wait() {
	<-l.done
}

func (l *Lexer) canceled() bool {
	select {
	case <-l.ctx.Done():
		return true
	default:
		return false
	}
}

func (l *Lexer) newToken(tokenType TokenType, ch byte) Token {
	// string([]byte{ch}) keeps an invalid UTF-8 byte as is; string(ch) would
	// turn it into the rune with that number.
	return Token{Type: tokenType, Literal: string([]byte{ch}), Line: l.Line, Column: l.Column, Pos: l.position}
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

	// Columns count characters: the trailing bytes of a valid multibyte
	// character stay in the column of its first byte, and each invalid
	// byte takes a column of its own.
	if l.trailingBytesInChar > 0 {
		l.trailingBytesInChar -= 1
		return
	}

	l.Column += 1
	l.currentCharIsInvalidUTF8 = false

	if l.currentChar >= utf8.RuneSelf {
		_, charSize := utf8.DecodeRuneInString(l.currentInput.content[l.position:])
		l.trailingBytesInChar = charSize - 1
		l.currentCharIsInvalidUTF8 = charSize == 1
	}
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.currentInput.content) {
		return 0
	} else {
		return l.currentInput.content[l.nextPosition]
	}
}

func (l *Lexer) eatWhitespace() {
	for isWhiteSpace(l.currentChar) && !l.canceled() {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	position := l.position
	for (isDigit(l.currentChar) || l.currentChar == '_') && !l.canceled() {
		l.readChar()
	}
	return l.currentInput.content[position:l.position]
}

func (l *Lexer) readRegexp() string {
	position := l.position + 1
	sentInvalidByte := false

	for {
		prevChar := l.currentChar
		l.readChar()
		if (l.currentChar == '/' && prevChar != '\\') || l.currentChar == 0 || l.canceled() {
			break
		}

		// Ruby rejects invalid UTF-8 in literals, but not in comments.
		if l.currentCharIsInvalidUTF8 && !sentInvalidByte {
			l.sendToken(l.newToken(ILLEGAL, l.currentChar))
			sentInvalidByte = true
		}
	}

	return l.currentInput.content[position:l.position]
}

func (l *Lexer) invalidUTF8At(pos int) bool {
	if pos >= len(l.currentInput.content) || l.currentInput.content[pos] < utf8.RuneSelf {
		return false
	}

	_, charSize := utf8.DecodeRuneInString(l.currentInput.content[pos:])
	return charSize == 1
}

func (l *Lexer) identifierStartsAt(pos int) bool {
	return identifierStartSize(l.currentInput.content, pos) > 0
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for !l.canceled() {
		charSize := identifierPartSize(l.currentInput.content, l.position)
		if charSize == 0 {
			break
		}

		for range charSize {
			l.readChar()
		}
	}

	if l.currentChar == '?' || l.currentChar == '!' {
		l.readChar()
	}

	return l.currentInput.content[position:l.position]
}

// FileName is the name of the input being lexed
func (l *Lexer) FileName() string {
	if l == nil || l.currentInput == nil {
		return ""
	}
	return l.currentInput.fileName
}
