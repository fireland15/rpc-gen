package schema_test

import (
	"testing"

	schema2 "github.com/fireland15/rpc-gen/internal/schema"
)

func TestNewMethodWithKindAndReturnType(t *testing.T) {
	userType := schema2.NamedType("User")
	arg1, _ := schema2.NewArgument("UserID", schema2.NamedType("string"))
	arg2, _ := schema2.NewArgument("UserName", schema2.NamedType("string"))

	tests := []struct {
		name       string
		methodName string
		kind       schema2.MethodKind
		returnType schema2.TypeRef
		args       []schema2.Argument
		wantErr    bool
		wantHTTP   schema2.HTTPMethod
	}{
		{
			name:       "Query with return type",
			methodName: "GetUser",
			kind:       schema2.MethodKindQuery,
			returnType: userType,
			args:       []schema2.Argument{arg1},
			wantErr:    false,
			wantHTTP:   schema2.HTTPGet,
		},
		{
			name:       "Mutation with return type",
			methodName: "CreateUser",
			kind:       schema2.MethodKindMutation,
			returnType: userType,
			args:       []schema2.Argument{arg1, arg2},
			wantErr:    false,
			wantHTTP:   schema2.HTTPPost,
		},
		{
			name:       "Stream without return type",
			methodName: "UserStream",
			kind:       schema2.MethodKindStream,
			returnType: nil,
			args:       []schema2.Argument{},
			wantErr:    false,
			wantHTTP:   schema2.HTTPPost,
		},
		{
			name:       "Empty method name",
			methodName: "",
			kind:       schema2.MethodKindQuery,
			returnType: nil,
			args:       []schema2.Argument{},
			wantErr:    true,
		},
		{
			name:       "Duplicate schema.Arguments",
			methodName: "DuplicateArgs",
			kind:       schema2.MethodKindMutation,
			returnType: nil,
			args:       []schema2.Argument{arg1, arg1},
			wantErr:    true,
		},
		{
			name:       "Nil schema.Argument",
			methodName: "NilArg",
			kind:       schema2.MethodKindQuery,
			returnType: nil,
			args:       []schema2.Argument{nil},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := schema2.NewMethod(tt.methodName, tt.kind, tt.returnType, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got=%v", tt.wantErr, err)
			}
			if tt.wantErr {
				return
			}

			if m.Name() != tt.methodName {
				t.Errorf("expected Name()=%s, got %s", tt.methodName, m.Name())
			}

			if m.Kind() != tt.kind {
				t.Errorf("expected Kind()=%v, got %v", tt.kind, m.Kind())
			}

			if m.HTTPMethod() != tt.wantHTTP {
				t.Errorf("expected schema.HTTPMethod()=%v, got %v", tt.wantHTTP, m.HTTPMethod())
			}

			if m.ReturnType() != tt.returnType {
				t.Errorf("expected ReturnType()=%v, got %v", tt.returnType, m.ReturnType())
			}

			argsMap := m.Arguments()
			if len(argsMap) != len(tt.args) {
				t.Errorf("expected %d schema.Arguments, got %d", len(tt.args), len(argsMap))
			}
		})
	}
}

func TestArgumentJSONNameAndType(t *testing.T) {
	// LowerCamelCase conversion
	arg, err := schema2.NewArgument("User_ID", schema2.NamedType("string"))
	if err != nil {
		t.Fatal(err)
	}
	if arg.JSONName() != "userId" {
		t.Errorf("expected JSONName userID, got %s", arg.JSONName())
	}

	nt, ok := arg.Type().(schema2.NamedTypeRef)
	if !ok {
		t.Errorf("expected NamedTypeRef, got %T", arg.Type())
	} else if nt.Name() != "string" {
		t.Errorf("expected Name()=string, got %s", nt.Name())
	}
}

func TestMethodArgumentsImmutability(t *testing.T) {
	arg1, _ := schema2.NewArgument("UserID", schema2.NamedType("string"))
	m, _ := schema2.NewMethod("TestMethod", schema2.MethodKindQuery, nil, arg1)

	args := m.Arguments()
	args["UserID"] = nil // mutate returned map

	if m.Arguments()["UserID"] == nil {
		t.Error("method schema.Arguments map should be immutable from outside")
	}
}
