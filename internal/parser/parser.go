package parser

import (
	"errors"
	"fmt"
	"io"
	"strings"

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

var ErrUnexpectedToken = errors.New("unexpected token")

type Keyword string

const (
	KwModel    Keyword = "model"
	KwOptional Keyword = "optional"
	KwRpc      Keyword = "rpc"
)

func (p *Parser) Parse() (model.ServiceDefinition, error) {
	def := model.ServiceDefinition{}
	parseErrors := make([]string, 0)

	for {
		tok, err := p.tokens.Lookahead(0)
		if err != nil {
			break
		}

		if tok.Text == string(KwModel) {
			md, err := p.parseModelDefinition()
			if err != nil {
				continue
			}
			def.Models = append(def.Models, md)
			continue
		} else if tok.Text == string(KwRpc) {
			rd, err := p.parseRpcDefinition()
			if err != nil {
				continue
			}
			def.Methods = append(def.Methods, rd)
			continue
		} else {
			msg := fmt.Sprintf("(%d:%d): expected keyword \"model\" or \"rpc\", but got \"%s\" instead", tok.Span.Start.Line, tok.Span.Start.Column, tok.Type)
			parseErrors = append(parseErrors, msg)
			p.tokens.Next()
		}
	}

	if len(parseErrors) > 0 {
		err := errors.New(strings.Join(parseErrors, "\n"))
		return def, err
	}
	return def, nil
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

	method.Name = rpcName

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
		method.Parameters = append(method.Parameters, parameter)

		tok, err = p.tokens.Lookahead(0)
		if err != nil || tok.Type != lexing.TokenTypeComma {
			break
		}
		p.tokens.Next()
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
	ty.Variant = model.TypeVariantNamed

	return p.parseOuterType(ty)
}

func (p *Parser) parseOuterType(inner model.Type) (model.Type, error) {
	tok, err := p.tokens.Lookahead(0)
	if err != nil {
		return inner, nil
	}

	if tok.Type == lexing.TokenTypeLeftSquareBracket {
		p.tokens.Next()
		if tok, err = p.tokens.Lookahead(0); err == nil {
			if tok.Type != lexing.TokenTypeRightSquareBracket {
				return inner, ErrUnexpectedToken
			}
			p.tokens.Next()
			new_type := model.Type{
				Variant: model.TypeVariantArray,
				Inner:   &inner,
			}
			return p.parseOuterType(new_type)
		}
	} else if tok.Type == lexing.TokenTypeQuestion {
		p.tokens.Next()
		new_type := model.Type{
			Name:    "",
			Variant: model.TypeVariantOptional,
			Inner:   &inner,
		}
		return p.parseOuterType(new_type)
	}
	return inner, nil
}

func (p *Parser) parseModelDefinition() (model.Model, error) {
	definition := model.Model{}
	err := p.parseKeyword(KwModel)
	if err != nil {
		return definition, err
	}

	modelName, err := p.parseIdentifier()
	if err != nil {
		return definition, err
	}

	definition.Name = modelName

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

func (p *Parser) parseModelFieldDefinition() (model.Field, error) {
	field := model.Field{}
	fieldName, err := p.parseIdentifier()
	if err != nil {
		return field, err
	}

	field.Name = fieldName

	fieldType, err := p.parseType()
	if err != nil {
		return field, err
	}

	field.Type = fieldType

	return field, nil
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

func (p *Parser) parseIdentifier() (string, error) {
	t, err := p.tokens.Next()
	if err != nil {
		return "", err
	}
	if t.Type != lexing.TokenTypeIdentifier {
		err = fmt.Errorf("expected identifier, but found \"%s\": %w", t.Text, ErrUnexpectedToken)
		return "", err
	}

	return t.Text, nil
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
