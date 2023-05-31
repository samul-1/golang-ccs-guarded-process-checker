package main

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Literals
	IDENT // identifiers: x, y, z, ...

	// Operators & punctuation
	PIPE          // | for composition
	PLUS          // + for summation
	BACKSLASH     // \ for restriction
	BRACKET_OPEN  // [ for relabeling
	BRACKET_CLOSE // ] for relabeling
	DOT

	LPAREN // (
	RPAREN // )

	// Keywords
	NIL
	REC
)

type Token struct {
	Type    TokenType
	Literal string
}

var eof = rune(0)

// utility function used to detect whitespace
func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// utility function used to read identifiers
func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// end of file detected
		l.ch = eof
	} else {
		l.ch = rune(l.input[l.readPosition])
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case eof:
		tok = Token{EOF, ""}
	case '|':
		tok = Token{PIPE, string(l.ch)}
	case '+':
		tok = Token{PLUS, string(l.ch)}
	case '\\':
		tok = Token{BACKSLASH, string(l.ch)}
	case '.':
		tok = Token{DOT, string(l.ch)}
	case '(':
		tok = Token{LPAREN, string(l.ch)}
	case ')':
		tok = Token{RPAREN, string(l.ch)}
	case '[':
		tok = Token{BRACKET_OPEN, string(l.ch)}
	case ']':
		tok = Token{BRACKET_CLOSE, string(l.ch)}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else {
			tok = Token{ILLEGAL, string(l.ch)}
		}
	}

	l.readChar()
	return tok
}

// utility function that consumes alphabetic chars to form an identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// utility function that advances as long as the current character is whitespace
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

// checks whether the given string is a keyword and return its token
// type if it is, otherwise returns token type IDENT
func lookupIdent(ident string) TokenType {
	switch ident {
	case "nil":
		return NIL
	case "rec":
		return REC
	default:
		return IDENT
	}
}
