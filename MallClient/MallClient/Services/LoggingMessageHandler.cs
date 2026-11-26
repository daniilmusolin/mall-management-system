using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using MallClient.Interfaces;

namespace MallClient.Services;

public class LoggingMessageHandler : DelegatingHandler {
    private readonly IAppLogger _logger;

    public LoggingMessageHandler(IAppLogger logger) {
        _logger = logger;
    }

    protected override async Task<HttpResponseMessage> SendAsync(HttpRequestMessage request, CancellationToken cancellationToken) {
        _logger.LogInfo($"HTTP {request.Method} {request.RequestUri}");

        try {
            var response = await base.SendAsync(request, cancellationToken);

            if (!response.IsSuccessStatusCode) {
                _logger.LogWarning($"HTTP {(int)response.StatusCode} {request.Method} {request.RequestUri}");
            }

            return response;
        } catch (HttpRequestException ex) {
            _logger.LogError($"HTTP Ошибка сети при запросе {request.RequestUri}", ex);
            throw;
        }
    }
}