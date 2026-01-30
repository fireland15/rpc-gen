using RpcGen.Examples.Endpoints.Abstractions;

namespace RpcGen.Examples.Services;

public class AuthService : ILoginHandler, IChangePasswordHandler
{
    Task ILoginHandler.HandleAsync(string password, string username, CancellationToken cancellationToken)
    {
        return Task.CompletedTask;
    }

    Task<bool> IChangePasswordHandler.HandleAsync(string newPassword, string oldPassword, CancellationToken cancellationToken)
    {
        return Task.FromResult(true);
    }
}