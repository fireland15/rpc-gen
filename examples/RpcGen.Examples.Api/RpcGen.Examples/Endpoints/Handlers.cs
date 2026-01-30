using Microsoft.AspNetCore.Mvc;
using RpcGen.Examples.Endpoints.Abstractions;
using RpcGen.Examples.Models;

namespace RpcGen.Examples.Endpoints;

public static class ProtocolEndpoints
{
    public static void MapProtocol(this WebApplication app)
    {
        app.MapPost("/login", HandleLogin);
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
