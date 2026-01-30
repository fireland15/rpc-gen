using Microsoft.AspNetCore.Mvc;
using RpcGen.Examples.Endpoints.Abstractions;
using RpcGen.Examples.Models;

namespace RpcGen.Examples.Endpoints;

public static class ProtocolEndpoints
{
    public static void MapProtocol(this WebApplication app)
    {
        app.MapPost("/change_password", HandleChangePassword);
        app.MapPost("/create_journal_entry", HandleCreateJournalEntry);
        app.MapPost("/extend_session", HandleExtendSession);
        app.MapGet("/get_journal_entry", HandleGetJournalEntry);
        app.MapPost("/signin", HandleSignin);
        app.MapPost("/signout", HandleSignout);
        app.MapPost("/upload_journal_entry_image", HandleUploadJournalEntryImage);
    }


    private static async Task<IResult> HandleChangePassword(
        [FromBody] ChangePasswordParameters parameters,
        [FromServices] IChangePasswordHandler handler,
        CancellationToken cancellationToken)
    {
        var result = await handler.HandleAsync(parameters.NewPassword, parameters.OldPassword, cancellationToken);
        return Results.Ok(result);
    }

    private static async Task<IResult> HandleCreateJournalEntry(
        [FromServices] ICreateJournalEntryHandler handler,
        CancellationToken cancellationToken)
    {
        var result = await handler.HandleAsync(cancellationToken);
        return Results.Ok(result);
    }

    private static async Task<IResult> HandleExtendSession(
        [FromServices] IExtendSessionHandler handler,
        CancellationToken cancellationToken)
    {
        await handler.HandleAsync(cancellationToken);
        return Results.NoContent();
    }

    private static async Task<IResult> HandleGetJournalEntry(
        [FromQuery(Name = "id")] Guid id,
        [FromServices] IGetJournalEntryHandler handler,
        CancellationToken cancellationToken)
    {
        var result = await handler.HandleAsync(id, cancellationToken);
        return Results.Ok(result);
    }

    private static async Task<IResult> HandleSignin(
        [FromBody] SigninParameters parameters,
        [FromServices] ISigninHandler handler,
        CancellationToken cancellationToken)
    {
        var result = await handler.HandleAsync(parameters.Password, parameters.Username, cancellationToken);
        return Results.Ok(result);
    }

    private static async Task<IResult> HandleSignout(
        [FromServices] ISignoutHandler handler,
        CancellationToken cancellationToken)
    {
        await handler.HandleAsync(cancellationToken);
        return Results.NoContent();
    }

    private static async Task<IResult> HandleUploadJournalEntryImage(
        [FromBody] UploadJournalEntryImageParameters parameters,
        [FromServices] IUploadJournalEntryImageHandler handler,
        CancellationToken cancellationToken)
    {
        await handler.HandleAsync(parameters.Image, parameters.Request, cancellationToken);
        return Results.NoContent();
    }

}
