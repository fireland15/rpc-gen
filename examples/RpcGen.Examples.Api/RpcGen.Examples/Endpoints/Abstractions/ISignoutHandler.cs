using RpcGen.Examples.Models;


namespace RpcGen.Examples.Endpoints.Abstractions;

public interface ISignoutHandler
{
    Task HandleAsync(CancellationToken cancellationToken = default);
}
