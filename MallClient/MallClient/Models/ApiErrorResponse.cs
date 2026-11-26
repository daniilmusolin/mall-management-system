using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace MallClient.Models;

public class ApiErrorResponse {
    [JsonPropertyName("error")]
    public string Error { get; set; } = string.Empty;

    [JsonPropertyName("fields")]
    public Dictionary<string, string>? Fields { get; set; }
}