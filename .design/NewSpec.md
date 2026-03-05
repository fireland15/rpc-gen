# New Idea for HTTP Rest IDL

## Example (WIP)

common.def
```
type ProblemDetail {
    Detail String
    Title String
    StatusCode Int
}

type Page<T> {
    Items [T!]!
    NextCursor String
}

responseset CommonResponses {
    // The request was not formatted correctly
    400 ProblemDetail!
        
    // The request was unauthenticated.
    401 empty
    
    // The Todo was not found.
    404 ProblemDetail!
        
    // There was a conflict with another resource.
    409 ProblemDetail!
        
    // The request could not be processed.
    422 ProblemDetail!
}
```

todos_service.def
```
import common

scalar Guid
scalar String
scalar Int

enum TodoStatus {
    PENDING
    DONE
}

type Todo {
    Id Guid!
    Status TodoStatus!
    Title String!
    Description String
}

type Project {
    Id Guid!
    Title String!
    Description String!
    Todos [Todo!]!
}

input CreateTodo {
    Title String!
    Description String
}

input CreateProject {
    Title String!
    Description String
}

input MoveTodoToProject {
    TodoId Guid!
    ProjectId Guid!
}

// Creates a todo item
@authorize
@version(1.23)
endpoint CreateTodo {
    POST /todos
    body CreateTodo!
    responses {
        // The request was successful
        201: Todo!
        
        ...common.CommonResponses
    }
}

@authorize
@version(1.25)
endpoint GetTodo {
    GET /todos/{todoId}
    path {
        // The ID of the Todo to get.
        todoId: Guid!
    }
    query {
        // A filter to apply to the todo
        filter String
    }
    responses {
        // The request was successful and the body contains the Todo.
        200: Todo!
        
        ...common.CommonResponses
    }
}

@authorize
@version(1.23)
endpoint ListTodos {
    GET /todos
    query {
        // A filter to apply to the Todos
        filter String
    }
    responses {
        // The request was successful and the body contains the Todo.
        200: common.Page<Todo>!
        
        ...common.CommonResponses
    }
}
```