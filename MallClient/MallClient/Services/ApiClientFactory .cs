using System.Net.Http;
using System.Net.Http.Headers;
using MallClient.Interfaces;

namespace MallClient.Services;

public class ApiClientFactory : IApiClientFactory {
    private const string ClientName = "MallClient";

    private readonly IHttpClientFactory _factory;

    public ApiClientFactory(IHttpClientFactory factory) {
        _factory = factory;
    }

    public HttpClient Create() {
        var client = _factory.CreateClient(ClientName);
        client.DefaultRequestHeaders.CacheControl =
            new CacheControlHeaderValue { NoCache = true, NoStore = true };
        return client;
    }
}