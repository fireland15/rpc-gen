package lexing

import (
	"errors"
	"fmt"
	"unicode"
)

type TokenType int

const (
	TokenTypeIdentifier TokenType = iota
	TokenTypeLeftBracket
	TokenTypeRightBracket
	TokenTypeLeftParenthesis
	TokenTypeRightParenthesis
	TokenTypeLeftSquareBracket
	TokenTypeRightSquareBracket
	TokenTypeQuestion
)

func (tt TokenType) String() string {
	switch tt {
	case TokenTypeIdentifier:
		return "identifier"
	case TokenTypeLeftBracket:
		return "{"
	case TokenTypeLeftParenthesis:
		return "("
	case TokenTypeLeftSquareBracket:
		return "["
	case TokenTypeRightBracket:
		return "}"
	case TokenTypeRightParenthesis:
		return ")"
	case TokenTypeRightSquareBracket:
		return "]"
	case TokenTypeQuestion:
		return "?"
	default:
		panic("unknown token type")
	}
}

var ErrInvalidToken = errors.New("invalid token")

type Token struct {
	Text string
	Span Span
	Type TokenType
}

type Span struct {
	Start Position
	End   Position
}

type Position struct {
	Offset int
	Line   int
	Column int
}

func NewToken(text string, span Span) (Token, error) {
	t := Token{
		Text: text,
		Span: span,
	}
	if isIdentifier(text) {
		t.Type = TokenTypeIdentifier
	} else if text == "{" {
		t.Type = TokenTypeLeftBracket
	} else if text == "}" {
		t.Type = TokenTypeRightBracket
	} else if text == "(" {
		t.Type = TokenTypeLeftParenthesis
	} else if text == ")" {
		t.Type = TokenTypeRightParenthesis
	} else if text == "[" {
		t.Type = TokenTypeLeftSquareBracket
	} else if text == "]" {
		t.Type = TokenTypeRightSquareBracket
	} else if text == "?" {
		t.Type = TokenTypeQuestion
	} else {
		err := fmt.Errorf("%s is not a recognized token: %w", text, ErrInvalidToken)
		return t, err
	}

	return t, nil
}

func isIdentifier(text string) bool {
	first := rune(text[0])
	return unicode.IsLetter(first) || first == '_'
}
