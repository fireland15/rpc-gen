using RpcGen.Examples.Models;


namespace RpcGen.Examples.Endpoints.Abstractions;

public interface IGetTodosHandler
{
    Task<Todo[]> HandleAsync(CancellationToken cancellationToken = default);
}
