# Design

**RPC Definition**

```
model Request {
    id: int
    name: string
    dueDate: date
    createdOn: datetime
}

rpc DoRequest(request Request) string
```

Gets transformed into...

**Protocol Definition**

```
{
    protocolDefinition: {
        types: [
            {
                name: "int",
                variant: "Scalar",
                serializedType: "number"
            },
            {
                name: "string",
                variant: "Scalar",
                serializedType: "string",
                serializedFormat: null
            },
            {
                name: "date",
                variant: "Scalar",
                serializedType: "string"
            },
            {
                name: "datetime",
                variant: "Scalar",
                serializedType: "string"
            },
            {
                name: "Request",
                variant: "Object",
                fields: [
                    {
                        name: "id",
                        type: "int"
                    },
                    {
                        name: "name",
                        type: "string"
                    },
                    {
                        name: "createdOn",
                        type: "date"
                    }
                ]
            },
            {
                name: "DoRequestParams",
                variant: "Object",
                fields: [
                    {
                        name: "request",
                        type: "Request"
                    }
                ]
            }
        ]
    }
}
```
