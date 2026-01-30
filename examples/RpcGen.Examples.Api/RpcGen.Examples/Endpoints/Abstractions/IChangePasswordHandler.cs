using RpcGen.Examples.Models;


namespace RpcGen.Examples.Endpoints.Abstractions;

public interface IChangePasswordHandler
{
    Task<ChangePasswordResponse> HandleAsync(string newPassword, string oldPassword, CancellationToken cancellationToken = default);
}
