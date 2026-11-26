using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public class TenantService : BaseApiService<Tenant>, ITenantService {
    public TenantService(IApiClientFactory clients, ISession session, IAppLogger logger)
        : base(clients, session, logger, "/api/tenants") { }
}