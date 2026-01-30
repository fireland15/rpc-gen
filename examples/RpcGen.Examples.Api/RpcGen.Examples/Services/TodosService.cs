using RpcGen.Examples.Endpoints.Abstractions;
using RpcGen.Examples.Models;

namespace RpcGen.Examples.Services;

public class TodosService : IGetTodosHandler, IGetTodoHandler
{
    public Task<Todo[]> HandleAsync(CancellationToken cancellationToken = default)
    {
        return Task.FromResult<Todo[]>([]);
    }

    public Task<Todo> HandleAsync(int id, CancellationToken cancellationToken = default)
    {
        return Task.FromResult(new Todo
        {
            Id = 23
        });
    }
}