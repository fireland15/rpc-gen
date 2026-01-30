using RpcGen.Examples.Models;


namespace RpcGen.Examples.Endpoints.Abstractions;

public interface ICreateJournalEntryHandler
{
    Task<JournalEntry> HandleAsync(CancellationToken cancellationToken = default);
}
