package main

import (
	"os"
	"strings"
	"text/template"

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
}

model CreateUserResponse {
	userId uuid
}
	
rpc GetUsers() UserPage
rpc CreateUser(User) CreateUserResponse
rpc PingUser(int)`

	c := compiler.Compiler{}
	service, err := c.Compile(strings.NewReader(source))
	if err != nil {
		panic(err)
	}

	templ, err := template.New("ts-api").ParseFiles("ts_client.template")
	if err != nil {
		panic(err)
	}
	err = templ.ExecuteTemplate(os.Stdout, "ts_client.template", service)
	if err != nil {
		panic(err)
	}
}
