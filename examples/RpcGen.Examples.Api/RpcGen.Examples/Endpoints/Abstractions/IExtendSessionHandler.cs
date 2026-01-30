using RpcGen.Examples.Models;


namespace RpcGen.Examples.Endpoints.Abstractions;

public interface IExtendSessionHandler
{
    Task HandleAsync(CancellationToken cancellationToken = default);
}
