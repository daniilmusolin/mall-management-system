using System.Net;
using System.Net.Http;
using System.Net.Http.Json;
using System.Text.Json;
using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public abstract class BaseApiService<T> : IApiCrudService<T> where T : class {
    protected readonly IApiClientFactory Clients;
    protected readonly string Endpoint;
    protected readonly ISession Session;
    private readonly IAppLogger _logger;

    protected BaseApiService(IApiClientFactory clients, ISession session, IAppLogger logger, string endpoint) {
        Clients = clients;
        Session = session;
        _logger = logger;
        Endpoint = endpoint;
    }

    protected async Task<TR> HandleResponse<TR>(HttpResponseMessage response) {
        if (response.StatusCode == HttpStatusCode.Unauthorized) {
            Session.RaiseSessionExpired();
            throw new ApiException("Требуется авторизация. Войдите снова.", 401);
        }

        if (!response.IsSuccessStatusCode) {
            var error = await response.Content.ReadAsStringAsync();
            try {
                var apiErr = JsonSerializer.Deserialize<ApiErrorResponse>(error);
                if (apiErr != null && !string.IsNullOrEmpty(apiErr.Error)) {
                    _logger.LogWarning($"API Ошибка ({response.StatusCode}): {apiErr.Error}");
                    throw new ApiException(apiErr.Error, (int)response.StatusCode);
                }
            } catch (JsonException) { }

            _logger.LogError($"API Ошибка ({response.StatusCode}): {error}");
            throw new ApiException($"Ошибка сервера: {(int)response.StatusCode}", (int)response.StatusCode);
        }

        if (typeof(TR) == typeof(EmptyResponse)) return (TR)(object)new EmptyResponse();
        if (response.StatusCode == HttpStatusCode.NoContent) return default!;

        var content = await response.Content.ReadAsStringAsync();
        if (string.IsNullOrWhiteSpace(content) || content.Trim() == "null") return default!;

        return JsonSerializer.Deserialize<TR>(content, new JsonSerializerOptions { PropertyNameCaseInsensitive = true })!;
    }

    public virtual Task<List<T>> GetAllAsync() => GetAllAsync(100, 0);

    public virtual async Task<List<T>> GetAllAsync(int limit = 100, int offset = 0) {
        var client = Clients.Create();
        var res = await client.GetAsync($"{Endpoint}?limit={limit}&offset={offset}");
        return await HandleResponse<List<T>>(res);
    }

    public virtual async Task CreateAsync(object data) {
        var client = Clients.Create();
        var res = await client.PostAsJsonAsync(Endpoint, data);
        await HandleResponse<EmptyResponse>(res);
    }

    public virtual async Task UpdateAsync(long id, object data) {
        var client = Clients.Create();
        var res = await client.PutAsJsonAsync($"{Endpoint}/{id}", data);
        await HandleResponse<EmptyResponse>(res);
    }

    public virtual async Task DeleteAsync(long id) {
        var client = Clients.Create();
        var res = await client.DeleteAsync($"{Endpoint}/{id}");
        await HandleResponse<EmptyResponse>(res);
    }
}

public class EmptyResponse { }