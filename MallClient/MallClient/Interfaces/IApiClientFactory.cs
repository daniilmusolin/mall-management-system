using System.Net.Http;

namespace MallClient.Interfaces;

public interface IApiClientFactory {
    HttpClient Create();
}