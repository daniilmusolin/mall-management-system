using MallClient.Models;

namespace MallClient.Interfaces;

public interface IAuthApi {
    Task<AuthResponse?> LoginAsync(string email, string password);
    Task<(bool Success, string ErrorBody)> RegisterAsync(string email, string password, string fullName);
    Task<AuthResponse?> RefreshAsync(string refreshToken);
}