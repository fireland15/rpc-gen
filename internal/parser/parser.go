package parser

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/fireland15/rpc-gen/internal/lexing"
	"github.com/fireland15/rpc-gen/internal/protocol"
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
	KwScalar   Keyword = "scalar"
	KwOptional Keyword = "optional"
	KwRpc      Keyword = "rpc"
)

func (p *Parser) Parse() (*protocol.Protocol, error) {
	def := protocol.NewProtocol("test")
	parseErrors := make([]string, 0)

	for {
		tok, err := p.tokens.Lookahead(0)
		if err != nil {
			break
		}

		if tok.Text == string(KwModel) {
			typeDef, err := p.parseModelDefinition()
			if err != nil {
				continue
			}
			def.AddTypeDefinition(typeDef)
			continue
		} else if tok.Text == string(KwScalar) {
			typedef, err := p.parseScalarDefinition()
			if err != nil {
				continue
			}
			if err := def.AddTypeDefinition(typedef); err != nil {
				slog.Error("parse error", slog.Any("err", err))
			}
			continue
		} else if tok.Text == string(KwRpc) {
			md, err := p.parseRpcDefinition()
			if err != nil {
				continue
			}
			if err := def.AddMethod(md); err != nil {
				slog.Error("parse error", slog.Any("err", err))
			}
			continue
		} else {
			msg := fmt.Sprintf("(%d:%d): expected keyword \"model\" or \"rpc\", but got \"%s\" instead", tok.Span.Start.Line, tok.Span.Start.Column, tok.Type)
			parseErrors = append(parseErrors, msg)
			p.tokens.Next()
		}
	}

	if len(parseErrors) > 0 {
		err := errors.New(strings.Join(parseErrors, "\n"))
		return nil, err
	}
	return def, nil
}

func (p *Parser) parseRpcDefinition() (*protocol.MethodDefinition, error) {
	err := p.parseKeyword(KwRpc)
	if err != nil {
		return nil, err
	}

	rpcName, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	method := protocol.NewMethodDefinition(rpcName)

	err = p.parseTokenType(lexing.TokenTypeLeftParenthesis)
	if err != nil {
		return method, err
	}

	for {
		tok, err := p.tokens.Lookahead(0)
		if err != nil || tok.Type != lexing.TokenTypeIdentifier {
			break
		}

		tok, err = p.tokens.Next()
		if err != nil {
			panic("lookahead failed?")
		}

		ty, err := p.parseTypeRef()
		if err != nil {
			return method, err
		}
		err = method.AddParameter(tok.Text, ty)
		if err != nil {
			slog.Error("parse error", slog.Any("err", err))
		}

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
			ty, err := p.parseTypeRef()
			if err != nil {
				return method, err
			}
			method.ReturnType = ty
		}
	}

	return method, nil
}

func (p *Parser) parseTypeRef() (*protocol.TypeRef2, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	ty := protocol.NewTypeRef(name)

	return p.parseOuterType(ty)
}

func (p *Parser) parseOuterType(inner *protocol.TypeRef2) (*protocol.TypeRef2, error) {
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
			new_type := protocol.NewArrayType(inner)
			return p.parseOuterType(new_type)
		}
	} else if tok.Type == lexing.TokenTypeQuestion {
		p.tokens.Next()
		new_type := protocol.NewOptionalType(inner)
		return p.parseOuterType(new_type)
	}
	return inner, nil
}

func (p *Parser) parseModelDefinition() (*protocol.TypeDefinition, error) {
	err := p.parseKeyword(KwModel)
	if err != nil {
		return nil, err
	}

	modelName, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	err = p.parseLeftBracket()
	if err != nil {
		return nil, err
	}

	td := protocol.NewObjectDefinition(modelName)

	for p.canParseModelFieldDefinition() {
		fd, err := p.parseModelFieldDefinition()
		if err != nil {
			break
		}

		err = td.AddField(fd)
		if err != nil {
			return td, err
		}
	}

	err = p.parseRightBracket()
	if err != nil {
		return td, err
	}

	return td, nil
}

func (p *Parser) parseScalarDefinition() (*protocol.TypeDefinition, error) {
	err := p.parseKeyword(KwScalar)
	if err != nil {
		return nil, err
	}

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	err = p.parseTokenType(lexing.TokenTypeEquals)
	if err != nil {
		return nil, err
	}

	serializedTypeStr, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	var serializedType protocol.SerializedType
	switch serializedTypeStr {
	case "number":
		serializedType = protocol.SerializedTypeNumber
	case "bool":
		serializedType = protocol.SerializedTypeBool
	case "string":
		serializedType = protocol.SerializedTypeString
	default:
		return nil, errors.New("unknown serialization type")
	}

	return protocol.NewScalarDefinition(name, serializedType), nil
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

func (p *Parser) parseModelFieldDefinition() (*protocol.FieldDefinition, error) {
	fieldName, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	fieldType, err := p.parseTypeRef()
	if err != nil {
		return nil, err
	}

	field := new(protocol.FieldDefinition)
	field.Name = fieldName
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
