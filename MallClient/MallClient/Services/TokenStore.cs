using MallClient.Interfaces;

namespace MallClient.Services;

public class TokenStore : ITokenStore {
    public string AccessToken { get; private set; } = string.Empty;
    public string RefreshToken { get; private set; } = string.Empty;

    public void SetTokens(string accessToken, string refreshToken) {
        AccessToken = accessToken;
        RefreshToken = refreshToken;
    }

    public void ClearTokens() {
        AccessToken = string.Empty;
        RefreshToken = string.Empty;
    }
}