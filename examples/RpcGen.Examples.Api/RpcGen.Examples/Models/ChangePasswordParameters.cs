using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models;

public record ChangePasswordParameters
{
    [JsonPropertyName("newPassword")]
    public string NewPassword { get; init; }
    [JsonPropertyName("oldPassword")]
    public string OldPassword { get; init; }
}
