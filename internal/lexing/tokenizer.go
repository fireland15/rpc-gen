package lexing

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

		if unicode.IsLetter(t.source.Current()) || t.source.Current() == '_' {
			text = append(text, t.source.Current())
			start := t.source.Position()
			for {
				err := t.source.Bump()
				if err != nil {
					return Token{
						Type: TokenTypeIdentifier,
						Text: string(text),
						Span: Span{
							Start: start,
							End:   t.source.Position(),
						},
					}, true
				}

				if !isIdentifierContinue(t.source.Current()) {
					return Token{
						Type: TokenTypeIdentifier,
						Text: string(text),
						Span: Span{
							Start: start,
							End:   t.source.Position(),
						},
					}, false
				}

				text = append(text, t.source.Current())
			}
		}

		if t.source.Current() == '{' {
			start := t.source.Position()
			err := t.source.Bump()
			return Token{
				Type: TokenTypeLeftBracket,
				Text: "{",
				Span: Span{
					Start: start,
					End:   start,
				},
			}, err != nil
		}

		if t.source.Current() == '}' {
			start := t.source.Position()
			err := t.source.Bump()
			return Token{
				Type: TokenTypeRightBracket,
				Text: "}",
				Span: Span{
					Start: start,
					End:   start,
				},
			}, err != nil
		}

		if t.source.Current() == '(' {
			start := t.source.Position()
			err := t.source.Bump()
			return Token{
				Type: TokenTypeLeftParenthesis,
				Text: "(",
				Span: Span{
					Start: start,
					End:   start,
				},
			}, err != nil
		}

		if t.source.Current() == ')' {
			start := t.source.Position()
			err := t.source.Bump()
			return Token{
				Type: TokenTypeRightParenthesis,
				Text: ")",
				Span: Span{
					Start: start,
					End:   start,
				},
			}, err != nil
		}

		if t.source.Current() == '[' {
			start := t.source.Position()
			err := t.source.Bump()
			return Token{
				Type: TokenTypeLeftSquareBracket,
				Text: "[",
				Span: Span{
					Start: start,
					End:   start,
				},
			}, err != nil
		}

		if t.source.Current() == ']' {
			start := t.source.Position()
			err := t.source.Bump()
			return Token{
				Type: TokenTypeRightSquareBracket,
				Text: "]",
				Span: Span{
					Start: start,
					End:   start,
				},
			}, err != nil
		}

		if t.source.Current() == '?' {
			start := t.source.Position()
			err := t.source.Bump()
			return Token{
				Type: TokenTypeQuestion,
				Text: "?",
				Span: Span{
					Start: start,
					End:   start,
				},
			}, err != nil
		}

		if t.source.Current() == ',' {
			start := t.source.Position()
			err := t.source.Bump()
			return Token{
				Type: TokenTypeComma,
				Text: ",",
				Span: Span{
					Start: start,
					End:   start,
				},
			}, err != nil
		}

		fmt.Printf("unrecognized character '%c'\n", t.source.Current())
		t.source.Bump()
	}
}

func isIdentifierContinue(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_'
}

type IRuneStream interface {
	Bump() error
	Current() rune
	Position() Position
}

type RuneStream struct {
	reader   *bufio.Reader
	current  rune
	position Position
}

func NewRuneStream(reader *bufio.Reader) (*RuneStream, error) {
	rs := new(RuneStream)
	rs.reader = reader
	rs.position = Position{
		Line:   0,
		Column: 0,
		Offset: 0,
	}
	rs.current = rune(0)
	err := rs.Bump()
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func (rs *RuneStream) Position() Position {
	return rs.position
}

func (rs *RuneStream) Current() rune {
	return rs.current
}

func (rs *RuneStream) Bump() error {
	r, _, err := rs.reader.ReadRune()
	if err != nil {
		return err
	}

	if rs.current != rune(0) {
		rs.position.Offset += 1
		rs.position.Column += 1
		if rs.current == '\n' {
			rs.position.Line += 1
			rs.position.Column = 0
		}
	}

	rs.current = r
	return nil
}
