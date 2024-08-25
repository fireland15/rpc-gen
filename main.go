package main

import (
	"fmt"
	"strings"

	"github.com/fireland15/rpc-gen/compiler"
)

func main() {
	source := `
model User {
	username string
	userId int
	email string optional
}
	
model UserPage {
	start int optional
	end int
	users User
}`

	c := compiler.Compiler{}
	service, err := c.Compile(strings.NewReader(source))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", service)
}
