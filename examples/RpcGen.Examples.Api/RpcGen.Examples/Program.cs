using RpcGen.Examples;
using RpcGen.Examples.Endpoints;
using RpcGen.Examples.Endpoints.Abstractions;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddScoped<IGetJournalEntryHandler, JournalService>();

var app = builder.Build();

app.MapProtocol();

app.Run();