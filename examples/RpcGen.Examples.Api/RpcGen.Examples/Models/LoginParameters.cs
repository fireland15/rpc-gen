using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record LoginParameters
{
    [JsonPropertyName("password")]
    public required string Password { get; init; }
    [JsonPropertyName("username")]
    public required string Username { get; init; }
}
