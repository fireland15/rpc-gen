package parser

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Tokenizer struct {
	source IRuneStream
}

func NewTokenizer(input io.Reader) (*Tokenizer, error) {
	r := bufio.NewReader(input)
	rs, err := NewRuneStream(r)
	if err != nil {
		return nil, err
	}
	return &Tokenizer{source: rs}, nil
}

// returns the next token
// the second value is true when at end of input
func (t *Tokenizer) Next() (Token, bool) {
	text := make([]rune, 0)
	for {
		// consume whitespace
		if unicode.Is(unicode.White_Space, t.source.Current()) {
			for {
				err := t.source.Bump()
				if err != nil {
					return Token{}, true
				}
				if !unicode.Is(unicode.White_Space, t.source.Current()) {
					break
				}
			}
		}

		if unicode.IsLetter(t.source.Current()) {
			text = append(text, t.source.Current())
			for {
				err := t.source.Bump()
				if err != nil {
					return Token{
						Type: TokenTypeIdentifier,
						Text: string(text),
					}, true
				}

				if !isIdentifierContinue(t.source.Current()) {
					return Token{
						Type: TokenTypeIdentifier,
						Text: string(text),
					}, false
				}

				text = append(text, t.source.Current())
			}
		}

		if t.source.Current() == '{' {
			err := t.source.Bump()
			return Token{
				Type: TokenTypeLeftBracket,
				Text: "{",
			}, err != nil
		}

		if t.source.Current() == '}' {
			err := t.source.Bump()
			return Token{
				Type: TokenTypeRightBracket,
				Text: "}",
			}, err != nil
		}

		if t.source.Current() == '(' {
			err := t.source.Bump()
			return Token{
				Type: TokenTypeLeftParenthesis,
				Text: "(",
			}, err != nil
		}

		if t.source.Current() == ')' {
			err := t.source.Bump()
			return Token{
				Type: TokenTypeRightParenthesis,
				Text: ")",
			}, err != nil
		}

		fmt.Printf("unrecognized character '%c'\n", t.source.Current())
		t.source.Bump()
	}
}

func isIdentifierContinue(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

type IRuneStream interface {
	Bump() error
	Current() rune
}

type RuneStream struct {
	reader  *bufio.Reader
	current rune
}

func NewRuneStream(reader *bufio.Reader) (*RuneStream, error) {
	rs := new(RuneStream)
	rs.reader = reader
	err := rs.Bump()
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (rs *RuneStream) Current() rune {
	return rs.current
}

func (rs *RuneStream) Bump() error {
	r, _, err := rs.reader.ReadRune()
	if err != nil {
		return err
	}
	rs.current = r
	return nil
}
