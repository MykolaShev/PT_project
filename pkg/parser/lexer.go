package parser

// Lexer performs lexical analysis, converting raw source code into a stream of tokens.
type Lexer struct {
	input        string
	position     int  // Current position in input (points to current char)
	readPosition int  // Current reading position in input (after current char)
	ch           byte // Current char under examination

	line   int // Current line number
	column int // Current column number
}

// NewLexer initializes a new Lexer with the given input string.
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances our position in the input string.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for "NUL" signifies EOF
	} else {
		l.ch = l.input[l.readPosition]
	}

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}

	l.position = l.readPosition
	l.readPosition++
}

// peekChar looks at the next character without advancing the lexer's position.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken evaluates the current character and returns the corresponding Token.
func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()

	startLine := l.line
	startColumn := l.column

	switch l.ch {
	case '=':
		tok = l.newToken(ASSIGN, l.ch) // For SIPL, '=' is equality. ':=' is assignment.
	case ';':
		tok = l.newToken(SEMICOLON, l.ch)
	case '(':
		tok = l.newToken(LPAREN, l.ch)
	case ')':
		tok = l.newToken(RPAREN, l.ch)
	case '+':
		tok = l.newToken(PLUS, l.ch)
	case '-':
		tok = l.newToken(MINUS, l.ch)
	case '*':
		tok = l.newToken(ASTERISK, l.ch)
	case '/':
		tok = l.newToken(SLASH, l.ch)
	case '<':
		tok = l.newToken(LT, l.ch)
	case '>':
		tok = l.newToken(GT, l.ch)
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: ASSIGN, Literal: string(ch) + string(l.ch)}
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
		}
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			tok.Line = startLine
			tok.Column = startColumn
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			tok.Line = startLine
			tok.Column = startColumn
			return tok
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
		}
	}

	tok.Line = startLine
	tok.Column = startColumn
	l.readChar()
	return tok
}

// skipWhitespace advances the lexer past any spaces, tabs, or newlines.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier reads a sequence of letters (a variable name or keyword).
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a sequence of digits (an integer literal).
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// newToken creates a Token struct with a single-character literal.
func (l *Lexer) newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

// Helpers for character classification
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}