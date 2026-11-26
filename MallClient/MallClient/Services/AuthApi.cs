using System.Net.Http.Json;
using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public class AuthApi : IAuthApi {
    private readonly IApiClientFactory _clients;

    public AuthApi(IApiClientFactory clients) {
        _clients = clients;
    }

    public async Task<AuthResponse?> LoginAsync(string email, string password) {
        var client = _clients.Create();
        var res = await client.PostAsJsonAsync("/api/auth/login", new { email, password });
        if (!res.IsSuccessStatusCode) return null;
        return await res.Content.ReadFromJsonAsync<AuthResponse>();
    }

    public async Task<(bool Success, string ErrorBody)> RegisterAsync(string email, string password, string fullName) {
        var client = _clients.Create();
        var res = await client.PostAsJsonAsync("/api/auth/register", new {
            email,
            password,
            full_name = fullName,
        });

        if (res.IsSuccessStatusCode)
            return (true, string.Empty);

        var body = await res.Content.ReadAsStringAsync();
        return (false, body);
    }

    public async Task<AuthResponse?> RefreshAsync(string refreshToken) {
        var client = _clients.Create();
        var res = await client.PostAsJsonAsync("/api/auth/refresh", new { refresh_token = refreshToken });
        if (!res.IsSuccessStatusCode) return null;
        return await res.Content.ReadFromJsonAsync<AuthResponse>();
    }
}