using RpcGen.Examples.Models;


namespace RpcGen.Examples.Endpoints.Abstractions;

public interface IChangePasswordHandler
{
    Task<bool> HandleAsync(string newPassword, string oldPassword, CancellationToken cancellationToken = default);
}
