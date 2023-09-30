package main

import (
	"fmt"
	"os"
	"unicode"
)

type TokenType int

const (
	TokenLParen TokenType = iota
	TokenRParen
	TokenComma
	TokenAnd
	TokenOr
	TokenIdent
	TokenString
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value string
}

func (t Token) IsLiteral() bool {
	return t.Type == TokenString || t.Type == TokenIdent
}

type Lexer struct {
	source string
	pos    int
}

func (l *Lexer) PeekChar() byte {
	if l.pos+1 >= len(l.source) {
		return 0
	}
	return l.source[l.pos+1]
}

func (l *Lexer) UnexpectedTokenError() {
	err := fmt.Errorf("[lex] unexpected token %b at position %d", l.source[l.pos], l.pos)
	fmt.Println(err)
	os.Exit(1)
}

func (l *Lexer) Next() Token {
	for l.pos < len(l.source) {
		ch := l.source[l.pos]

		if unicode.IsSpace(rune(ch)) {
			l.pos++
			continue
		}

		switch ch {
		case '(':
			l.pos++
			return Token{Type: TokenLParen, Value: string(ch)}
		case ')':
			l.pos++
			return Token{Type: TokenRParen, Value: string(ch)}
		case ',':
			l.pos++
			return Token{Type: TokenComma, Value: string(ch)}
		case '&':
			if l.PeekChar() == '&' {
				l.pos += 2
				return Token{Type: TokenAnd, Value: "&&"}
			} else {
				l.UnexpectedTokenError()
			}
		case '|':
			if l.PeekChar() == '|' {
				l.pos += 2
				return Token{Type: TokenOr, Value: "||"}
			} else {
				l.UnexpectedTokenError()
			}
		default:
			if ch == '"' || ch == '\'' {
				quoteChar := ch
				start := l.pos
				l.pos++
				for l.pos < len(l.source) && l.source[l.pos] != quoteChar {
					l.pos++
				}
				l.pos++
				// capture only value inside quote
				return Token{Type: TokenString, Value: l.source[start+1 : l.pos-1]}
			} else {
				start := l.pos
				for l.pos < len(l.source) && !unicode.IsSpace(rune(l.source[l.pos])) {
					l.pos++
				}
				return Token{Type: TokenIdent, Value: l.source[start:l.pos]}
			}
		}
	}
	return Token{Type: TokenEOF, Value: ""}
}

func (l *Lexer) LexAll() []Token {
	var tokens []Token
	for {
		token := l.Next()
		tokens = append(tokens, token)
		if token.Type == TokenEOF {
			break
		}
	}
	return tokens
}
