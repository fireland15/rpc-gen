using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record UploadJournalEntryImageRequest
{
    [JsonPropertyName("journalEntryId")]
    public int JournalEntryId { get; init; }
}
