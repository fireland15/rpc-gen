using RpcGen.Examples.Models;


namespace RpcGen.Examples.Endpoints.Abstractions;

public interface IGetJournalEntryHandler
{
    Task<JournalEntry> HandleAsync(Guid id, CancellationToken cancellationToken = default);
}
