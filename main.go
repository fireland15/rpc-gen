package main

import (
	"os"
	"strings"
	"text/template"

	"github.com/fireland15/rpc-gen/compiler"
	"github.com/fireland15/rpc-gen/writer"
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

	cfgas, err := writer.ReadConfig("config.json")
	if err != nil {
		panic(err)
	}

	print(cfgas)

	templ, err := template.New("ts-api").ParseFiles("ts_client.template")
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile("service_client.ts", os.O_CREATE, 0)
	if err != nil {
		panic(err)
	}

	err = templ.ExecuteTemplate(f, "ts_client.template", service)
	if err != nil {
		panic(err)
	}
}
