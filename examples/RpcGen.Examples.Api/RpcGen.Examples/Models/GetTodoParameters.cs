using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record GetTodoParameters
{
    [JsonPropertyName("id")]
    public required int Id { get; init; }
}
