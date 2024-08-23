package compiler

type Type interface {
	Name() string
}

type Rpc struct {
	Name         string
	RequestType  *Type
	ResponseType *Type
}

type Service struct {
	rpc   map[string]Rpc
	types map[string]Type
}

func Compile() (*Service, error) {

}

type Object struct {
	name string
}

type Scalar struct {
	name string
}
