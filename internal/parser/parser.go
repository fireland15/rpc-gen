package parser

import (
	"errors"
	"fmt"
	"io"
)

type Parser struct {
	tokens *TokenStream
}

func NewParser(input io.Reader) (*Parser, error) {
	ts, err := NewTokenStream(input)
	if err != nil {
		return nil, err
	}

	p := new(Parser)
	p.tokens = ts

	return p, nil
}

type ServiceDefinition struct {
	Models []ModelDefinition
	Rpc    []RpcDefinition
}

type ModelDefinition struct {
	Name   string
	Fields []FieldDefinition
}

type FieldDefinition struct {
	Name     string
	TypeName string
	Optional bool
}

type RpcDefinition struct {
	Name             string
	RequestTypeName  string
	ResponseTypeName string
}

var ErrUnexpectedToken = errors.New("unexpected token")

type Keyword string

const (
	KwModel    Keyword = "model"
	KwOptional Keyword = "optional"
	KwRpc      Keyword = "rpc"
)

func (p *Parser) Parse() (ServiceDefinition, error) {
	def := ServiceDefinition{}
	var err error

	for {
		tok, err := p.tokens.Lookahead(0)
		if err != nil {
			break
		}

		if tok.Type != TokenTypeIdentifier {
			p.tokens.Next()
			continue
		} else if tok.Text == string(KwModel) {
			md, err := p.parseModelDefinition()
			if err != nil {
				break
			}
			def.Models = append(def.Models, md)
			continue
		} else if tok.Text == string(KwRpc) {
			rd, err := p.parseRpcDefinition()
			if err != nil {
				break
			}
			def.Rpc = append(def.Rpc, rd)
			continue
		} else {
			p.tokens.Next()
		}
	}

	return def, err
}

func (p *Parser) parseRpcDefinition() (RpcDefinition, error) {
	def := RpcDefinition{}
	err := p.parseKeyword(KwRpc)
	if err != nil {
		return def, err
	}

	rpcName, err := p.parseIdentifier()
	if err != nil {
		return def, err
	}

	def.Name = rpcName.Text

	err = p.parseTokenType(TokenTypeLeftParenthesis)
	if err != nil {
		return def, err
	}

	n, err := p.tokens.Lookahead(0)
	if err == nil {
		if n.Type == TokenTypeIdentifier {
			def.RequestTypeName = n.Text
			p.tokens.Next()
		}
	}

	err = p.parseTokenType(TokenTypeRightParenthesis)
	if err != nil {
		return def, err
	}

	n, err = p.tokens.Lookahead(0)
	if err == nil {
		if n.Type == TokenTypeIdentifier && n.Text != string(KwModel) && n.Text != string(KwRpc) && n.Text != string(KwOptional) {
			def.ResponseTypeName = n.Text
			p.tokens.Next()
		}
	}

	return def, nil
}

func (p *Parser) parseModelDefinition() (ModelDefinition, error) {
	definition := ModelDefinition{}
	err := p.parseKeyword(KwModel)
	if err != nil {
		return definition, err
	}

	modelName, err := p.parseIdentifier()
	if err != nil {
		return definition, err
	}

	definition.Name = modelName.Text

	err = p.parseLeftBracket()
	if err != nil {
		return definition, err
	}

	for p.canParseModelFieldDefinition() {
		fd, err := p.parseModelFieldDefinition()
		if err != nil {
			break
		}

		definition.Fields = append(definition.Fields, fd)
	}

	err = p.parseRightBracket()
	if err != nil {
		return definition, err
	}

	return definition, nil
}

func (p *Parser) canParseModelFieldDefinition() bool {
	t, err := p.tokens.Lookahead(0)
	if err != nil {
		return false
	}

	if t.Type != TokenTypeIdentifier {
		return false
	}

	return true
}

func (p *Parser) parseModelFieldDefinition() (FieldDefinition, error) {
	definition := FieldDefinition{}
	t, err := p.parseIdentifier()
	if err != nil {
		return definition, err
	}

	definition.Name = t.Text

	t, err = p.parseIdentifier()
	if err != nil {
		return definition, err
	}

	definition.TypeName = t.Text

	t, err = p.tokens.Lookahead(0)
	if err != nil {
		return definition, nil
	}

	if t.Type != TokenTypeIdentifier || t.Text != string(KwOptional) {
		return definition, nil
	}

	t, err = p.tokens.Next()
	if err != nil {
		return definition, err
	}

	definition.Optional = true

	return definition, nil
}

func (p *Parser) parseLeftBracket() error {
	return p.parseTokenType(TokenTypeLeftBracket)
}

func (p *Parser) parseRightBracket() error {
	return p.parseTokenType(TokenTypeRightBracket)
}

func (p *Parser) parseTokenType(tt TokenType) error {
	t, err := p.tokens.Next()
	if err != nil {
		return err
	}

	if t.Type != tt {
		err = fmt.Errorf("expected \"%s\", but found \"%s\": %w", tt.String(), t.Text, ErrUnexpectedToken)
		return err
	}

	return nil
}

func (p *Parser) parseIdentifier() (Token, error) {
	t, err := p.tokens.Next()
	if err != nil {
		return Token{}, err
	}
	if t.Type != TokenTypeIdentifier {
		err = fmt.Errorf("expected identifier, but found \"%s\": %w", t.Text, ErrUnexpectedToken)
		return Token{}, err
	}
	return t, nil
}

func (p *Parser) parseKeyword(kw Keyword) error {
	t, err := p.tokens.Next()
	if err != nil {
		return err
	}
	if t.Type != TokenTypeIdentifier {
		err = fmt.Errorf("expected keyword \"%s\", but found \"%s\": %w", kw, t.Text, ErrUnexpectedToken)
		return err
	}
	if t.Text != string(kw) {
		err = fmt.Errorf("expected keyword \"%s\", but got \"%s\": %w", kw, t.Text, ErrUnexpectedToken)
		return err
	}
	return nil
}
