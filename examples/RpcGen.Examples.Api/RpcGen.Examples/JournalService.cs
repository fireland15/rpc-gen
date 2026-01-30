using RpcGen.Examples.Endpoints.Abstractions;
using RpcGen.Examples.Models;

namespace RpcGen.Examples;

public class JournalService : IGetJournalEntryHandler
{
    public Task<JournalEntry> HandleAsync(Guid id, CancellationToken cancellationToken = default)
    {
        return Task.FromResult(new JournalEntry
        {
            CreatedOn = DateOnly.FromDateTime(DateTime.UtcNow),
            Details = "Here are some details about your thing.",
            Id = Guid.NewGuid(),
            Status = 123,
            Title = "The journal entry is about things that suck",
            UpdatedOn = default
        });
    }
}