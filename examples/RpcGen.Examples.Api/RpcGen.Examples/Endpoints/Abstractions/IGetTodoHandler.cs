using RpcGen.Examples.Models;


namespace RpcGen.Examples.Endpoints.Abstractions;

public interface IGetTodoHandler
{
    Task<Todo> HandleAsync(int id, CancellationToken cancellationToken = default);
}
