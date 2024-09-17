package lexing

import (
	"errors"
	"fmt"
	"io"
)

var ErrEndOfStream = errors.New("end of token stream")

type TokenStream struct {
	tokenizer *Tokenizer
	end       bool
	buf       *RingBuffer[Token]
}

func NewTokenStream(input io.Reader) (*TokenStream, error) {
	tokenizer, err := NewTokenizer(input)
	if err != nil {
		return nil, err
	}

	ts := new(TokenStream)
	ts.buf = NewRingBuffer[Token](10)
	ts.end = false
	ts.tokenizer = tokenizer
	return ts, nil
}

func (t *TokenStream) Next() (Token, error) {
	if !t.buf.IsEmpty() {
		tok, err := t.buf.Pop()
		if err != nil {
			err = fmt.Errorf("buffer is empty after checking it was not: %w", err)
			panic(err)
		}
		return tok, nil
	}

	if t.end {
		return Token{}, ErrEndOfStream
	}

	tok, end := t.tokenizer.Next()
	if end {
		t.end = true
	}
	return tok, nil
}

func (t *TokenStream) Lookahead(n int) (Token, error) {
	for t.buf.Size() <= n {
		if t.end {
			return Token{}, ErrEndOfStream
		}

		tok, end := t.tokenizer.Next()
		if end {
			t.end = true
		}
		t.buf.Push(tok)
	}

	tok, ok := t.buf.At(n)
	if !ok {
		panic("shouldn't be here")
	}

	return tok, nil
}
