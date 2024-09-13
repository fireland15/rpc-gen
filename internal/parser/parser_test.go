package parser

import (
	"strings"
	"testing"
)

func TestParserParsesModelDefinition(t *testing.T) {
	source := `
model Banans {
	name string
	stuff int optional
}`
	p, err := NewParser(strings.NewReader(source))
	if err != nil {
		t.Error(err)
	}

	md, err := p.parseModelDefinition()
	if err != nil {
		t.Error(err)
	}

	ExpectEqual(t, "model definition name", "Banans", md.Name)
	ExpectEqual(t, "model definition field count", 2, len(md.Fields))

	f := md.Fields[0]
	if f.Name != "name" || f.TypeName != "string" || f.Optional != false {
		t.Error("field not parsed correctly")
	}

	f = md.Fields[1]
	if f.Name != "stuff" || f.TypeName != "int" || f.Optional != true {
		t.Error("field not parsed correctly")
	}
}

func TestParserParsesRpcDefinition(t *testing.T) {
	source := `
rpc Do(Soemthing) void`

	p, err := NewParser(strings.NewReader(source))
	if err != nil {
		t.Error(err)
	}

	def, err := p.parseRpcDefinition()
	if err != nil {
		t.Error(err)
	}

	if def.Name != "Do" {
		t.Error("rpc definition name not parsed correctly")
	}

	if def.RequestTypeName != "Soemthing" {
		t.Error("rpc parameter type name not parsed correctly")
	}

	if def.ResponseTypeName != "void" {
		t.Error("rpc response type name not parsed correctly")
	}
}

func TestParserParsesFullRpcFile(t *testing.T) {
	source := `
model SignInRequest {
	username string
	password string
}

model SignInResponse {
	token string
	expires instant
}
	
rpc SignIn(SignInRequest) SignInResponse

rpc SignOut() void`

	p, err := NewParser(strings.NewReader(source))
	if err != nil {
		t.Error(err)
	}

	def, err := p.Parse()
	if err != nil {
		t.Error(err)
		return
	}

	ExpectEqual(t, "model count", 2, len(def.Models))
	ExpectEqual(t, "rpc count", 2, len(def.Rpc))
}
