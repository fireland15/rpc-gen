using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record JournalEntry
{
    [JsonPropertyName("createdOn")]
    public DateOnly CreatedOn { get; init; }
    [JsonPropertyName("details")]
    public string? Details { get; init; }
    [JsonPropertyName("id")]
    public Guid Id { get; init; }
    [JsonPropertyName("status")]
    public int Status { get; init; }
    [JsonPropertyName("title")]
    public string Title { get; init; }
    [JsonPropertyName("updatedOn")]
    public DateOnly? UpdatedOn { get; init; }
}
