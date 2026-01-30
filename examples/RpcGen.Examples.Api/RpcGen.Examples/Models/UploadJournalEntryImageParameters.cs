using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record UploadJournalEntryImageParameters
{
    [JsonPropertyName("image")]
    public int Image { get; init; }
    [JsonPropertyName("request")]
    public UploadJournalEntryImageRequest Request { get; init; }
}
