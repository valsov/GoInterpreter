package lexer

import (
	"errors"
	"strings"

	"github.com/valsov/gointerpreter/token"
)

type Lexer struct {
	input        []rune
	position     int  // Position in input pointing to current char (ch)
	readPosition int  // Reading position in input (peeker)
	ch           rune // Current char
}

func New(input string) *Lexer {
	lexer := &Lexer{input: []rune(input)}
	lexer.readChar()
	return lexer
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpaces()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok.Type = token.EQ
			tok.Literal = "=="
			l.readChar() // Read again to advance past the peeked char
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			tok.Type = token.NOT_EQ
			tok.Literal = "!="
			l.readChar()
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	case '"':
		str, err := l.readString()
		if err != nil {
			tok.Type = token.ILLEGAL
			tok.Literal = err.Error()
		} else {
			tok.Type = token.STRING
			tok.Literal = str
		}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readWhile(isLetter)
			tok.Type = token.Lookup(tok.Literal)
			// readIdentifier() already advanced read pointers, no need to call readChar() -> return early
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readWhile(isDigit)
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar() // Call readChar to advance to next character for subsequent uses
	return tok
}

// Read next char from the input, loading it in the lexer and avdancing the read pointers
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) skipWhiteSpaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readWhile(filterFunc func(rune) bool) string {
	position := l.position
	for filterFunc(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position]) // l.position is used instead of l.readPosition because we are already pointing to the next (invalid) char
}

func (l *Lexer) readString() (string, error) {
	sb := strings.Builder{}
	for {
		l.readChar()

		switch l.ch {
		case 0:
			return "", errors.New("encountered EOF before string end")
		case '\\':
			l.readChar()
			switch l.ch {
			case '\\':
				sb.WriteByte('\\')
			case '"':
				sb.WriteByte('"')
			case 't':
				sb.WriteByte('\t')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			default:
				sb.WriteByte('\\')
				sb.WriteRune(l.ch)
			}
		case '"':
			return sb.String(), nil
		default:
			sb.WriteRune(l.ch)
		}
	}
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}
