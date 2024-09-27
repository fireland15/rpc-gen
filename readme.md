# RPC-Gen

Generate client and server code from a common definition file.

## Supported Languages/Frameworks

More will be added as needed.

| Language   | Framework | Server | Client |
| ---------- | --------- | ------ | ------ |
| Typescript | N/A       | ❌     | ✅     |
| Go         | Echo      | ✅     | ❌     |

## Definitions

Definition files describe the interface between clients and servers. RPC-Gen translates this definition into interfaces & implementations in various languages for both clients and servers.

The primary part of a definition file is a RPC statement. It looks something like:

```
rpc AssignUser(request AssignUserRequest) AssignUserResponse
```

The `rpc` keyword indicates the statement is for an RPC method. Then, the name of the method, its parameters, and finally the return value. The return value can be omitted.

The definition file also defines the data models of the service. For example:

```
model AssignUserRequest {
    userId    uuid
    projectId uuid
    makeOwner bool?
    tags      string[]
}
```

Models can have one or more fields with scalar or model types. Fields can be marked optional with the `optional` keyword.

### Built-in Scalar Types

| Type   | Go        | TS      |
| ------ | --------- | ------- |
| bool   | bool      | boolean |
| int    | int       | number  |
| float  | float     | number  |
| string | string    | string  |
| date   | time.Time | string  |

Custom types can be defined in the config file for each language.

## Usage

`go run ./cmd/cli/main.go -c ./config.json`
