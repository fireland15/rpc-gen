using RpcGen.Examples.Models;


namespace RpcGen.Examples.Endpoints.Abstractions;

public interface ILoginHandler
{
    Task HandleAsync(string password, string username, CancellationToken cancellationToken = default);
}
