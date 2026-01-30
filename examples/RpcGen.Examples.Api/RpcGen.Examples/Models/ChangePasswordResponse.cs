using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record ChangePasswordResponse
{
    [JsonPropertyName("date")]
    public DateOnly Date { get; init; }
    [JsonPropertyName("details")]
    public string Details { get; init; }
    [JsonPropertyName("name")]
    public string Name { get; init; }
}
