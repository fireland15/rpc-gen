using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record SigninResponse
{
    [JsonPropertyName("errors")]
    public string[] Errors { get; init; }
}
