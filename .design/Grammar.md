# Grammar

## File

```
File
  ::= ( Import
      | ScalarDecl
      | EnumDecl
      | TypeDecl
      | InputDecl
      | ResponseSetDecl
      | EndpointDecl
      | Annotation
      )* EOF
```

## Imports and Namespacing

Imports resolve relative to the current file. So in the example `import common`, the definition file at `./common.def` will be imported. All symbols imported from a file are namespaced to the file's name. So in our example, all symbols imported are available under the `common` namespace.

```
Import
::= "import" Identifier ";"

QualifiedName
::= Identifier ( "." Identifier )*
```

example:
```
import common;
common.Page<Todo>
```

## Identifiers & literals

```
Identifier
  ::= Letter ( Letter | Digit | "_" )*

IntegerLiteral
  ::= Digit+

StringLiteral
  ::= '"' ( ~["\\] | EscapeSequence )* '"'
  
FloatLiteral
  ::= Digit+ "." Digit+

Literal
  ::= StringLiteral
     | IntegerLiteral
     | FloatLiteral
     | "true"
     | "false"
```

## Annotations

```
Annotation
  ::= "@" Identifier AnnotationArgs?

AnnotationArgs
  ::= "(" AnnotationArg ( "," AnnotationArg )* ")"

AnnotationArg
  ::= Identifier "=" Literal
     | Literal
```

examples:
```
@authorize
@version(1.23)
@rateLimit(limit=10, unit="s")
```

## Scalars

```
ScalarDecl
  ::= "scalar" Identifier ScalarBody?

ScalarBody
  ::= "{" Annotation* "}"
```

examples:
```
scalar Guid
scalar Guid { @format("uuid") }
```

## Enums

```
EnumDecl
  ::= "enum" Identifier "{" EnumValue* "}"

EnumValue
  ::= Identifier
```

examples:
```
enum TodoStatus {
    PENDING
    DONE
}
```

## Types & Inputs
```
TypeDecl
  ::= "type" Identifier TypeParams? "{" Field* "}"

InputDecl
  ::= "input" Identifier TypeParams? "{" Field* "}"

TypeParams
  ::= "<" Identifier ( "," Identifier )* ">"
  
Field
  ::= Identifier TypeRef FieldNullability? Annotation*
  
TypeRef
  ::= QualifiedName TypeArgs?
    | ArrayTypeRef

TypeArgs
  ::= "<" TypeRef ( "," TypeRef )* ">"
  
FieldNullability
  ::= "!"
  
ArrayTypeRef
  ::= "[" TypeRef "]"
```

## Response Sets

```
ResponseSetDecl
  ::= "responseset" Identifier "{" ResponseEntry* "}"
  
ResponseEntry
  ::= StatusCode ":" ResponseBody
  
StatusCode
  ::= IntegerLiteral
  
ResponseBody
  ::= TypeRef "!"
     | "empty"
```

examples:
```
responseset CommonResponses {
    400: ProblemDetail!
    401: empty
}
```

## Endpoints

```
EndpointDecl
  ::= Annotation*
      "endpoint" Identifier
      "{"
          HttpMethod PathTemplate
          EndpointSection*
      "}"

HttpMethod
  ::= "GET" | "POST" | "PUT" | "PATCH" | "DELETE"

PathTemplate
  ::= "/" PathSegment ( "/" PathSegment )*
  
PathSegment
  ::= Identifier
     | "{" Identifier "}"
     
EndpointSection
  ::= PathParams
     | QueryParams
     | BodyDecl
     | ResponsesDecl

PathParams
  ::= "path" "{" ParamField* "}"

QueryParams
  ::= "query" "{" ParamField* "}"

ParamField
  ::= Identifier ":" TypeRef FieldNullability? Annotation*

BodyDecl
  ::= "body" TypeRef "!"

ResponsesDecl
  ::= "responses" "{"
        ( ResponseEntry
        | ResponseSpread
        )*
      "}"

ResponseSpread
  ::= "..." QualifiedName
```

## Whitespace

```
Comment
  ::= "//" ~[\n]*

Whitespace
  ::= ( " " | "\t" | "\n" | "\r" | Comment )+
```