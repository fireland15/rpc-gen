package parser

import (
	"errors"
	"fmt"
	"io"

	"github.com/fireland15/rpc-gen/internal/ast"
	"github.com/fireland15/rpc-gen/internal/lexing"
	"github.com/fireland15/rpc-gen/internal/model"
)

type Parser struct {
	tokens *lexing.TokenStream
}

func NewParser(input io.Reader) (*Parser, error) {
	ts, err := lexing.NewTokenStream(input)
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

		if tok.Type != lexing.TokenTypeIdentifier {
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

func (p *Parser) parseRpcDefinition() (model.Method, error) {
	method := model.Method{}
	err := p.parseKeyword(KwRpc)
	if err != nil {
		return method, err
	}

	rpcName, err := p.parseIdentifier()
	if err != nil {
		return method, err
	}

	method.Name = rpcName.Text

	err = p.parseTokenType(lexing.TokenTypeLeftParenthesis)
	if err != nil {
		return method, err
	}

	for {
		parameter := model.MethodParameter{}
		tok, err := p.tokens.Lookahead(0)
		if err != nil || tok.Type != lexing.TokenTypeIdentifier {
			break
		}

		tok, err = p.tokens.Next()
		if err != nil {
			panic("lookahead failed?")
		}
		parameter.Name = tok.Text

		ty, err := p.parseType()
		if err != nil {
			return method, err
		}
		parameter.Type = ty
	}

	err = p.parseTokenType(lexing.TokenTypeRightParenthesis)
	if err != nil {
		return method, err
	}

	next, err := p.tokens.Lookahead(0)
	if err == nil {
		if next.Type == lexing.TokenTypeIdentifier && !isKeyword(next.Text) {
			ty, err := p.parseType()
			if err != nil {
				return method, err
			}
			method.ReturnType = &ty
		}
	}

	return method, nil
}

func (p *Parser) parseType() (model.Type, error) {
	ty := model.Type{}

	name, err := p.parseIdentifier()
	if err != nil {
		return ty, err
	}

	ty.Name = name
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

	if t.Type != lexing.TokenTypeIdentifier {
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

	if t.Type != lexing.TokenTypeIdentifier || t.Text != string(KwOptional) {
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
	return p.parseTokenType(lexing.TokenTypeLeftBracket)
}

func (p *Parser) parseRightBracket() error {
	return p.parseTokenType(lexing.TokenTypeRightBracket)
}

func (p *Parser) parseTokenType(tt lexing.TokenType) error {
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

func (p *Parser) parseIdentifier() (ast.Identifier, error) {
	ident := ast.Identifier{}

	t, err := p.tokens.Next()
	if err != nil {
		return ident, err
	}
	if t.Type != lexing.TokenTypeIdentifier {
		err = fmt.Errorf("expected identifier, but found \"%s\": %w", t.Text, ErrUnexpectedToken)
		return ident, err
	}

	ident.Name = t.Text
	ident.Span = t.Span

	return ident, nil
}

func (p *Parser) parseKeyword(kw Keyword) error {
	t, err := p.tokens.Next()
	if err != nil {
		return err
	}
	if t.Type != lexing.TokenTypeIdentifier {
		err = fmt.Errorf("expected keyword \"%s\", but found \"%s\": %w", kw, t.Text, ErrUnexpectedToken)
		return err
	}
	if t.Text != string(kw) {
		err = fmt.Errorf("expected keyword \"%s\", but got \"%s\": %w", kw, t.Text, ErrUnexpectedToken)
		return err
	}
	return nil
}

func isKeyword(str string) bool {
	return str == string(KwModel) || str == string(KwRpc) || str == string(KwOptional)
}
