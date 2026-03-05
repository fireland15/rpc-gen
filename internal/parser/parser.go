package parser

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/fireland15/rpc-gen/internal/lexing"
	"github.com/fireland15/rpc-gen/internal/schema"
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
	KwEnum     Keyword = "enum"
	KwOptional Keyword = "optional"
	KwRpc      Keyword = "rpc"
	KwStream   Keyword = "stream"
)

func (p *Parser) Parse() (schema.Schema, error) {
	def := schema.NewSchema()
	parseErrors := make([]string, 0)

	for {
		tok, err := p.tokens.Lookahead(0)
		if err != nil {
			break
		}

		if tok.Text == string(KwModel) {
			modelDefinition, err := p.parseModelDefinition()
			if err != nil {
				continue
			}
			def.AddType(modelDefinition)
			continue
		} else if tok.Text == string(KwScalar) {
			scalarDefinition, err := p.parseScalarDefinition()
			if err != nil {
				slog.Error(err.Error())
				continue
			}
			if err := def.AddType(scalarDefinition); err != nil {
				slog.Error("parse error", slog.Any("err", err))
			}
			continue
		} else if tok.Text == string(KwEnum) {
			scalarDefinition, err := p.parseEnumDefinition()
			if err != nil {
				continue
			}
			if err := def.AddType(scalarDefinition); err != nil {
				slog.Error("parse error", slog.Any("err", err))
			}
			continue
		} else if tok.Text == string(KwRpc) {
			md, err := p.parseRpcDefinition(KwRpc)
			if err != nil {
				continue
			}
			if err := def.AddMethod(md); err != nil {
				slog.Error("parse error", slog.Any("err", err))
			}
			continue
		} else if tok.Text == string(KwStream) {
			md, err := p.parseRpcDefinition(KwStream)
			if err != nil {
				continue
			}
			if err := def.AddMethod(md); err != nil {
				slog.Error("parse error", slog.Any("err", err))
			}
			continue
		} else if tok.Type == lexing.TokenTypeAt {
			decorators, err := p.parseDecorators()
			if err != nil {
				continue
			}

			tok, err := p.tokens.Lookahead(0)
			if err != nil {
				continue
			}

			var md schema.Method
			switch tok.Text {
			case string(KwRpc):
				md, err = p.parseRpcDefinition(KwRpc)
			case string(KwStream):
				md, err = p.parseRpcDefinition(KwStream)
			default:
				parseErrors = append(parseErrors,
					fmt.Sprintf("(%d:%d): decorator not attached to method",
						tok.Span.Start.Line, tok.Span.Start.Column))
				continue
			}

			if err != nil {
				continue
			}

			for _, d := range decorators {
				md.AddDecorator(d)
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
	return def.Schema(), nil
}

func (p *Parser) parseRpcDefinition(kw Keyword) (schema.Method, error) {
	err := p.parseKeyword(kw)
	if err != nil {
		return nil, err
	}

	rpcName, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	err = p.parseTokenType(lexing.TokenTypeLeftParenthesis)
	if err != nil {
		return nil, err
	}

	args := make([]schema.Argument, 0)

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
			return nil, err
		}

		arg, err := schema.NewArgument(tok.Text, ty)
		if err != nil {
			return nil, err
		}

		args = append(args, arg)

		tok, err = p.tokens.Lookahead(0)
		if err != nil || tok.Type != lexing.TokenTypeComma {
			break
		}
		p.tokens.Next()
	}

	err = p.parseTokenType(lexing.TokenTypeRightParenthesis)
	if err != nil {
		return nil, err
	}

	var returnType schema.TypeRef
	next, err := p.tokens.Lookahead(0)
	if err == nil {
		if next.Type == lexing.TokenTypeIdentifier && !isKeyword(next.Text) {
			ty, err := p.parseTypeRef()
			if err != nil {
				return nil, err
			}
			returnType = ty
		}
	}

	method, err := schema.NewMethod(rpcName, toMethodKind(kw), returnType, args...)
	if err != nil {
		return nil, err
	}

	return method, nil
}

func toMethodKind(kw Keyword) schema.MethodKind {
	switch kw {
	case KwRpc:
		return schema.MethodKindRpc
	case KwStream:
		return schema.MethodKindStream
	default:
		panic("unsupported keyword")
	}
}

func (p *Parser) parseTypeRef() (schema.TypeRef, error) {
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	ty := schema.NamedType(name)

	return p.parseOuterType(ty)
}

func (p *Parser) parseOuterType(inner schema.TypeRef) (schema.TypeRef, error) {
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
			newType := schema.Array(inner)
			return p.parseOuterType(newType)
		}
	} else if tok.Type == lexing.TokenTypeQuestion {
		p.tokens.Next()
		newType := schema.Optional(inner)
		return p.parseOuterType(newType)
	}
	return inner, nil
}

func (p *Parser) parseModelDefinition() (schema.Composite, error) {
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

	td, err := schema.NewComposite(modelName)
	if err != nil {
		return nil, err
	}

	for p.canParseModelFieldDefinition() {
		fd, err := p.parseModelFieldDefinition()
		if err != nil {
			break
		}

		err = td.AddField(fd)
		if err != nil {
			return nil, err
		}
	}

	err = p.parseRightBracket()
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (p *Parser) parseScalarDefinition() (schema.Scalar, error) {
	err := p.parseKeyword(KwScalar)
	if err != nil {
		return nil, err
	}

	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	scalar := schema.NewScalar(name)

	return scalar, nil
}

func (p *Parser) parseEnumDefinition() (schema.Enumeration, error) {
	err := p.parseKeyword(KwEnum)
	if err != nil {
		return nil, err
	}
	enumName, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err = p.parseLeftBracket(); err != nil {
		return nil, err
	}

	var variants = make([]string, 0)

	for p.canParseEnumVariant() {
		variantName, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		variants = append(variants, variantName)
	}

	if err = p.parseRightBracket(); err != nil {
		return nil, err
	}

	enum, err := schema.NewEnumeration(enumName, variants...)
	if err != nil {
		return nil, err
	}

	return enum, nil
}

func (p *Parser) canParseEnumVariant() bool {
	t, err := p.tokens.Lookahead(0)
	if err != nil {
		return false
	}

	if t.Type != lexing.TokenTypeIdentifier {
		return false
	}
	return true
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

func (p *Parser) parseModelFieldDefinition() (schema.Field, error) {
	fieldName, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	fieldType, err := p.parseTypeRef()
	if err != nil {
		return nil, err
	}

	field := schema.NewField(fieldName, fieldType)

	return field, nil
}

func (p *Parser) parseDecorators() ([]schema.Decorator, error) {
	decorators := make([]schema.Decorator, 0)

	for {
		tok, err := p.tokens.Lookahead(0)
		if err != nil || tok.Type != lexing.TokenTypeAt {
			break
		}

		// consume '@'
		p.tokens.Next()

		name, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}

		args := make([]string, 0)

		// optional argument list
		if tok, _ := p.tokens.Lookahead(0); tok.Type == lexing.TokenTypeLeftParenthesis {
			p.tokens.Next() // '('

			for {
				arg, err := p.parseIdentifier()
				if err != nil {
					return nil, err
				}
				args = append(args, arg)

				tok, err := p.tokens.Lookahead(0)
				if err != nil || tok.Type != lexing.TokenTypeComma {
					break
				}
				p.tokens.Next() // ','
			}

			if err := p.parseTokenType(lexing.TokenTypeRightParenthesis); err != nil {
				return nil, err
			}
		}

		d, err := schema.NewDecorator(name, args)
		if err != nil {
			return nil, err
		}
		decorators = append(decorators, d)
	}

	return decorators, nil
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
	return str == string(KwModel) ||
		str == string(KwRpc) ||
		str == string(KwOptional) ||
		str == string(KwStream) ||
		str == string(KwEnum)
}
