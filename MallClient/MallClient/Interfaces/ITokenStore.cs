namespace MallClient.Interfaces;

public interface ITokenStore {
    string AccessToken { get; }
    string RefreshToken { get; }
    void SetTokens(string accessToken, string refreshToken);
    void ClearTokens();
}