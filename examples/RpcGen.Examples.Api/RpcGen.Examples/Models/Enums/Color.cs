using System.Text.Json.Serialization;

namespace RpcGen.Examples.Models.Enums;

public enum Color
{
    [JsonStringEnumMemberName("RED")]
    Red,
    [JsonStringEnumMemberName("BLUE")]
    Blue,
    [JsonStringEnumMemberName("GREEN")]
    Green,
}
