using System.Threading.Tasks;
using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public class AuthService : IAuthService {
    private readonly IAuthApi _api;
    private readonly IApiErrorParser _errorParser;
    private readonly ISession _session;
    private readonly ITokenStore _tokenStore;

    public AuthService(IAuthApi api, IApiErrorParser errorParser, ISession session, ITokenStore tokenStore) {
        _api = api;
        _errorParser = errorParser;
        _session = session;
        _tokenStore = tokenStore;
    }

    public async Task<bool> LoginAsync(string email, string password) {
        var response = await _api.LoginAsync(email, password);
        if (response == null) return false;

        _tokenStore.SetTokens(response.Token, response.RefreshToken);
        _session.Set(response.Token, response.User);
        return true;
    }

    public async Task<(bool Success, string Error)> RegisterAsync(string email, string password, string fullName) {
        var (success, errorBody) = await _api.RegisterAsync(email, password, fullName);
        if (success) return (true, string.Empty);
        return (false, _errorParser.Parse(errorBody));
    }

    public void Logout() {
        _tokenStore.ClearTokens();
        _session.Clear();
    }
}