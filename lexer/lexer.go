package lexer

import (
	"fmt"

	"github.com/dustinlindquist/monkey/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	char         byte
}

func New(input string) *Lexer {
	var l = &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
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

func (l *Lexer) readNumber() string {
	var start = l.position
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) readIdentifier() string {
	var start = l.position
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\r' || l.char == '\n' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			var char = l.char
			l.readChar()
			var literal = string(char) + string(char)
			tok = token.Token{token.EQ, literal}
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '!':
		if l.peekChar() == '=' {
			var char = l.char
			l.readChar()
			var literal = string(char) + string(l.char)
			tok = token.Token{token.NOT_EQ, literal}
		} else {
			tok = newToken(token.BANG, l.char)
		}
	case '/':
		tok = newToken(token.SLASH, l.char)
	case '<':
		tok = newToken(token.LT, l.char)
	case '>':
		tok = newToken(token.GT, l.char)
	case '*':
		tok = newToken(token.ASTERISK, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			fmt.Println("expected a space")
			tok = newToken(token.ILLEGAL, l.char)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{tokenType, string(char)}
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
