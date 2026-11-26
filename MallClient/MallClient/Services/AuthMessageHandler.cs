using System.Net;
using System.Net.Http;
using System.Net.Http.Headers;
using System.Threading;
using System.Threading.Tasks;
using MallClient.Interfaces;

namespace MallClient.Services;

public class AuthMessageHandler : DelegatingHandler {
    private static readonly SemaphoreSlim RefreshLock = new(1, 1);

    private readonly ITokenStore _tokenStore;
    private readonly ISession _session;
    private readonly IAuthApi _authApi;

    public AuthMessageHandler(ITokenStore tokenStore, ISession session, IAuthApi authApi) {
        _tokenStore = tokenStore;
        _session = session;
        _authApi = authApi;
    }

    protected override async Task<HttpResponseMessage> SendAsync(HttpRequestMessage request, CancellationToken cancellationToken) {
        var tokenBefore = _tokenStore.AccessToken;
        if (!string.IsNullOrEmpty(tokenBefore)) {
            request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", tokenBefore);
        }

        var response = await base.SendAsync(request, cancellationToken);

        if (response.StatusCode == HttpStatusCode.Unauthorized &&
            !request.RequestUri!.AbsolutePath.Contains("/api/auth/login") &&
            !request.RequestUri!.AbsolutePath.Contains("/api/auth/refresh")) {

            bool refreshed;

            await RefreshLock.WaitAsync(cancellationToken);
            try {
                if (_tokenStore.AccessToken != tokenBefore) {
                    refreshed = true;
                } else {
                    refreshed = await TryRefreshTokenAsync(cancellationToken);
                }
            } finally {
                RefreshLock.Release();
            }

            if (!refreshed) {
                _tokenStore.ClearTokens();
                _session.RaiseSessionExpired();
                return response;
            }

            request.Headers.Authorization = new AuthenticationHeaderValue("Bearer", _tokenStore.AccessToken);
            return await base.SendAsync(request, cancellationToken);
        }

        return response;
    }

    private async Task<bool> TryRefreshTokenAsync(CancellationToken cancellationToken) {
        if (string.IsNullOrEmpty(_tokenStore.RefreshToken))
            return false;

        try {
            var authResponse = await _authApi.RefreshAsync(_tokenStore.RefreshToken);
            if (authResponse == null) return false;

            _tokenStore.SetTokens(authResponse.Token, authResponse.RefreshToken);
            _session.Set(authResponse.Token, _session.CurrentUser);
            return true;
        } catch {
            return false;
        }
    }
}