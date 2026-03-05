using Microsoft.AspNetCore.Mvc;
using RpcGen.Examples.Endpoints.Abstractions;
using RpcGen.Examples.Models;

namespace RpcGen.Examples.Endpoints;

public static class ProtocolEndpoints
{
    public static void MapProtocol(this WebApplication app)
    {
        app.MapPost("/change_password", HandleChangePassword);
        app.MapPost("/get_todo", HandleGetTodo);
        app.MapPost("/get_todos", HandleGetTodos);
        app.MapPost("/login", HandleLogin);
    }


    private static async Task<IResult> HandleChangePassword(
        [FromBody] ChangePasswordParameters parameters,
        [FromServices] IChangePasswordHandler handler,
        CancellationToken cancellationToken)
    {
        var result = await handler.HandleAsync(parameters.NewPassword, parameters.OldPassword, cancellationToken);
        return Results.Ok(result);
    }

    private static async Task<IResult> HandleGetTodo(
        [FromBody] GetTodoParameters parameters,
        [FromServices] IGetTodoHandler handler,
        CancellationToken cancellationToken)
    {
        var result = await handler.HandleAsync(parameters.Id, cancellationToken);
        return Results.Ok(result);
    }

    private static async Task<IResult> HandleGetTodos(
        [FromServices] IGetTodosHandler handler,
        CancellationToken cancellationToken)
    {
        var result = await handler.HandleAsync(cancellationToken);
        return Results.Ok(result);
    }

    private static async Task<IResult> HandleLogin(
        [FromBody] LoginParameters parameters,
        [FromServices] ILoginHandler handler,
        CancellationToken cancellationToken)
    {
        await handler.HandleAsync(parameters.Password, parameters.Username, cancellationToken);
        return Results.NoContent();
    }

}
