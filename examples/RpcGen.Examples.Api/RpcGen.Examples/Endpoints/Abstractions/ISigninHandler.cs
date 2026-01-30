using RpcGen.Examples.Models;


namespace RpcGen.Examples.Endpoints.Abstractions;

public interface ISigninHandler
{
    Task<SigninResponse> HandleAsync(string password, string username, CancellationToken cancellationToken = default);
}
