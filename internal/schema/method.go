package schema

import (
	"errors"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

type MethodKind string

const (
	MethodKindRpc    MethodKind = "rpc"
	MethodKindStream MethodKind = "stream"
)

type HTTPMethod string

const (
	HTTPGet  HTTPMethod = "GET"
	HTTPPost HTTPMethod = "POST"
)

////////////////////////////////
// Method
////////////////////////////////

type Method interface {
	// Name returns the logical name of the Method
	Name() string

	// Arguments returns the named arguments of the Method
	Arguments() map[string]Argument

	// Kind returns the method's MethodKind
	Kind() MethodKind

	// Path returns the URL path for the method call
	Path() string

	// HTTPMethod returns the HTTP method based on the Method's MethodKind
	HTTPMethod() HTTPMethod

	// ReturnType returns the TypeRef for the methods return. Returns nil when there is not returned value.
	ReturnType() TypeRef

	// Decorators returns a slice of Decorator instances attached to the method
	Decorators() []Decorator

	// AddDecorator adds a new decorator to the method
	AddDecorator(Decorator)
}

type method struct {
	name       string
	arguments  map[string]*argument
	kind       MethodKind
	returnType TypeRef
	decorators []Decorator
}

var _ Method = (*method)(nil)

// NewMethod constructs a new Method with the given name and arguments.
// Returns an error if name is empty or if there are duplicate argument names.
func NewMethod(name string, kind MethodKind, returnType TypeRef, arguments ...Argument) (Method, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("method name cannot be empty")
	}

	argMap := make(map[string]*argument)
	for _, arg := range arguments {
		if arg == nil {
			return nil, errors.New("argument cannot be nil")
		}

		if _, exists := argMap[arg.Name()]; exists {
			return nil, errors.New("duplicate argument name: " + arg.Name())
		}
		argMap[arg.Name()] = &argument{
			name:     arg.Name(),
			jsonName: arg.JSONName(),
			ty:       arg.Type(),
		}
	}

	return &method{
		name:       name,
		arguments:  argMap,
		kind:       kind,
		returnType: returnType,
		decorators: make([]Decorator, 0),
	}, nil
}

// Name implements Method
func (m *method) Name() string {
	return m.name
}

// Arguments implements Method
func (m *method) Arguments() map[string]Argument {
	// convert internal map of concrete arguments to map of interface
	result := make(map[string]Argument, len(m.arguments))
	for k, v := range m.arguments {
		result[k] = v
	}
	return result
}

func (m *method) Path() string {
	return fmt.Sprintf("/%s", strcase.ToSnake(m.name))
}

// Kind implements Method
func (m *method) Kind() MethodKind {
	return m.kind
}

// HTTPMethod implements Method
func (m *method) HTTPMethod() HTTPMethod {
	switch m.kind {
	case MethodKindRpc:
		return HTTPPost
	case MethodKindStream:
		return HTTPPost
	default:
		panic("unknown MethodKind: " + string(m.kind))
	}
}

// ReturnType implements Method
func (m *method) ReturnType() TypeRef {
	return m.returnType
}

func (m *method) Decorators() []Decorator {
	return m.decorators
}

func (m *method) AddDecorator(d Decorator) {
	m.decorators = append(m.decorators, d)
}

////////////////////////////////
// Argument
////////////////////////////////

type Argument interface {
	// Name returns the logical name of the Argument
	Name() string

	// JSONName is the name of the argument as sent in JSON over the wire.
	JSONName() string

	// Type returns the TypeRef for the type of the Argument
	Type() TypeRef
}

type argument struct {
	name     string
	jsonName string
	ty       TypeRef
}

var _ Argument = (*argument)(nil)

// NewArgument constructs a new Argument with the given name and TypeRef.
// Automatically converts the name to lowerCamelCase for JSON.
func NewArgument(name string, ty TypeRef) (Argument, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("argument name cannot be empty")
	}
	if ty == nil {
		return nil, errors.New("argument type cannot be nil")
	}

	return &argument{
		name:     name,
		jsonName: strcase.ToLowerCamel(name),
		ty:       ty,
	}, nil
}

// Name implements Argument
func (a *argument) Name() string {
	return a.name
}

// JSONName implements Argument
func (a *argument) JSONName() string {
	return a.jsonName
}

// Type implements Argument
func (a *argument) Type() TypeRef {
	return a.ty
}

////////////////////////////////
// Decorator
////////////////////////////////

type Decorator interface {
	// Name returns the logical name of the Decorator.
	Name() string

	// Arguments returns the comma separated values passed to the decorator.
	Arguments() []string
}

type decorator struct {
	name      string
	arguments []string
}

var _ Decorator = (*decorator)(nil)

func NewDecorator(name string, arguments []string) (Decorator, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("decorator name cannot be empty")
	}
	for idx, arg := range arguments {
		if arg == "" {
			return nil, fmt.Errorf("argument %d, cannot be empty", idx)
		}
	}
	return &decorator{
		name:      name,
		arguments: arguments,
	}, nil
}

func (d decorator) Name() string {
	return d.name
}

func (d decorator) Arguments() []string {
	return d.arguments
}
