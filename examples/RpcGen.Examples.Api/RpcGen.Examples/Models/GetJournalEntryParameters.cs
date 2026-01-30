using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record GetJournalEntryParameters
{
    [JsonPropertyName("id")]
    public Guid Id { get; init; }
}
