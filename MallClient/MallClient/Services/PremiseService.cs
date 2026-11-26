using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public class PremiseService : BaseApiService<Premise>, IPremiseService {
    public PremiseService(IApiClientFactory clients, ISession session, IAppLogger logger)
        : base(clients, session, logger, "/api/premises") { }
}