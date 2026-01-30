using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record ChangePasswordParameters
{
    [JsonPropertyName("newPassword")]
    public required string NewPassword { get; init; }
    [JsonPropertyName("oldPassword")]
    public required string OldPassword { get; init; }
}
