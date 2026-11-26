using System.Text.Json.Serialization;

namespace MallClient.Models;

public class AuthResponse {
    [JsonPropertyName("token")] public string Token { get; set; } = string.Empty;
    [JsonPropertyName("refresh_token")] public string RefreshToken { get; set; } = string.Empty;
    [JsonPropertyName("user")] public User User { get; set; } = new();
}

public class User {
    [JsonPropertyName("id")] public long Id { get; set; }
    [JsonPropertyName("email")] public string Email { get; set; } = string.Empty;
    [JsonPropertyName("full_name")] public string FullName { get; set; } = string.Empty;
    [JsonPropertyName("role")] public string Role { get; set; } = string.Empty;
}