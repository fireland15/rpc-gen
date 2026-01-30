using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record Todo
{
    [JsonPropertyName("id")]
    public required int Id { get; init; }
}
