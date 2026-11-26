using Microsoft.Extensions.DependencyInjection;
using MallClient.Interfaces;

namespace MallClient.Services;

public class CrudServiceFactory : ICrudServiceFactory {
    private readonly IServiceProvider _provider;

    public CrudServiceFactory(IServiceProvider provider) {
        _provider = provider;
    }

    public IApiCrudService<T> GetService<T>() where T : class {
        return _provider.GetRequiredService<IApiCrudService<T>>();
    }
}