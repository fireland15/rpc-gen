package parser

import (
	"errors"
	"strings"
	"testing"
)

func TestTokenStreamNextReturnsEndofStreamErr(t *testing.T) {
	ts, err := NewTokenStream(strings.NewReader("apples"))
	if err != nil {
		t.Error(err)
		return
	}
	_, _ = ts.Next()
	_, err = ts.Next()
	if !errors.Is(err, ErrEndOfStream) {
		t.Error(err)
	}
}

func TestTokenStreamNextReturnsToken(t *testing.T) {
	ts, err := NewTokenStream(strings.NewReader("something"))
	if err != nil {
		t.Error(err)
		return
	}
	tok, err := ts.Next()
	if err != nil {
		t.Error(err)
	}

	ExpectEqual(t, "token text", "something", tok.Text)
	ExpectEqual(t, "token type", TokenTypeIdentifier, tok.Type)
}

func TestTokenStreamLookaheadReturnsErrEndOfStream(t *testing.T) {
	ts, err := NewTokenStream(strings.NewReader("apples"))
	if err != nil {
		t.Error(err)
		return
	}
	_, err = ts.Lookahead(3)
	if !errors.Is(err, ErrEndOfStream) {
		t.Error(err)
	}
}

func TestTokenStreamLookaheadReturnsToken(t *testing.T) {
	ts, err := NewTokenStream(strings.NewReader("something something2"))
	if err != nil {
		t.Error(err)
	}

	tok, err := ts.Lookahead(1)
	if err != nil {
		t.Error(err)
	}

	ExpectEqual(t, "token text", "something2", tok.Text)
	ExpectEqual(t, "token type", TokenTypeIdentifier, tok.Type)
}
