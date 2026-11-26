using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public class RequestService : BaseApiService<ServiceRequest>, IRequestService {
    public RequestService(IApiClientFactory clients, ISession session, IAppLogger logger)
        : base(clients, session, logger, "/api/requests") { }
}