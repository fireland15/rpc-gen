package lexing

import (
	"strings"
	"testing"
)

func TestTokenizerTokenizes(t *testing.T) {
	tokenizer, err := NewTokenizer(strings.NewReader("apples { () }"))
	if err != nil {
		t.Error(err)
	}

	tok, end := tokenizer.Next()
	ExpectEqual(t, "end", false, end)
	ExpectEqual(t, "token type", TokenTypeIdentifier, tok.Type)
	ExpectEqual(t, "token text", "apples", tok.Text)

	tok, end = tokenizer.Next()
	ExpectEqual(t, "end", false, end)
	ExpectEqual(t, "token type", TokenTypeLeftBracket, tok.Type)
	ExpectEqual(t, "token text", "{", tok.Text)

	tok, end = tokenizer.Next()
	ExpectEqual(t, "end", false, end)
	ExpectEqual(t, "token type", TokenTypeLeftParenthesis, tok.Type)
	ExpectEqual(t, "token text", "(", tok.Text)

	tok, end = tokenizer.Next()
	ExpectEqual(t, "end", false, end)
	ExpectEqual(t, "token type", TokenTypeRightParenthesis, tok.Type)
	ExpectEqual(t, "token text", ")", tok.Text)

	tok, end = tokenizer.Next()
	ExpectEqual(t, "end", true, end)
	ExpectEqual(t, "token type", TokenTypeRightBracket, tok.Type)
	ExpectEqual(t, "token text", "}", tok.Text)
}

func TestTokenizerTokenizesRequiresDecorator(t *testing.T) {
	tokenizer, err := NewTokenizer(strings.NewReader(`@requires(policy)`))
	if err != nil {
		t.Error(err)
	}

	tok, end := tokenizer.Next()
	ExpectEqual(t, "end", false, end)
	ExpectEqual(t, "token type", TokenTypeAt, tok.Type)
	ExpectEqual(t, "token text", "@", tok.Text)

	tok, end = tokenizer.Next()
	ExpectEqual(t, "end", false, end)
	ExpectEqual(t, "token type", TokenTypeIdentifier, tok.Type)
	ExpectEqual(t, "token text", "requires", tok.Text)

	tok, end = tokenizer.Next()
	ExpectEqual(t, "end", false, end)
	ExpectEqual(t, "token type", TokenTypeLeftParenthesis, tok.Type)
	ExpectEqual(t, "token text", "(", tok.Text)

	tok, end = tokenizer.Next()
	ExpectEqual(t, "end", false, end)
	ExpectEqual(t, "token type", TokenTypeIdentifier, tok.Type)
	ExpectEqual(t, "token text", "policy", tok.Text)

	tok, end = tokenizer.Next()
	ExpectEqual(t, "end", true, end)
	ExpectEqual(t, "token type", TokenTypeRightParenthesis, tok.Type)
	ExpectEqual(t, "token text", ")", tok.Text)
}
