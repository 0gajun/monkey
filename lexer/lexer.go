package lexer

import "github.com/0gajun/monkey/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		lit, err := l.readString()
		if err != nil {
			tok = newToken(token.ILLEGAL, l.ch)
		} else {
			tok.Literal = lit
			tok.Type = token.STRING
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	return l.readWhile(isLetter)
}

func (l *Lexer) readNumber() string {
	return l.readWhile(isDigit)
}

type IllegalEscapedCharacterError byte

func (e IllegalEscapedCharacterError) Error() string {
	return "hogehoge"
}

func (l *Lexer) readString() (string, error) {
	l.readChar()

	var bytes = make([]byte, 0, 1000)

	for l.ch != '"' {
		var c byte
		if l.ch == '\\' {
			escapedChar, err := l.readEscapedCharacter()
			if err != nil {
				return "", err
			}
			c = escapedChar
		} else {
			c = l.ch
		}

		bytes = append(bytes, c)
		l.readChar()
	}

	return string(bytes), nil
}

func (l *Lexer) readEscapedCharacter() (byte, error) {
	l.readChar()

	var char byte
	switch l.ch {
	case '"':
		char = '"'
	case 't':
		char = '\t'
	case 'n':
		char = '\n'
	case '\\':
		char = '\\'
	default:
		return ' ', IllegalEscapedCharacterError(l.ch)
	}
	return char, nil
}

func (l *Lexer) readWhile(isContinue func(ch byte) bool) string {
	start_pos := l.position
	for isContinue(l.ch) {
		l.readChar()
	}
	return l.input[start_pos:l.position]
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *Lexer) Input() string {
	return l.input
}
