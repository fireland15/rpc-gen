using RpcGen.Examples.Endpoints;
using RpcGen.Examples.Endpoints.Abstractions;
using RpcGen.Examples.Services;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddScoped<ILoginHandler, AuthService>();
builder.Services.AddScoped<IChangePasswordHandler, AuthService>();
builder.Services.AddScoped<IGetTodosHandler, TodosService>();
builder.Services.AddScoped<IGetTodoHandler, TodosService>();

var app = builder.Build();

app.MapProtocol();

app.Run();