using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record SigninParameters
{
    [JsonPropertyName("password")]
    public string Password { get; init; }
    [JsonPropertyName("username")]
    public string Username { get; init; }
}
