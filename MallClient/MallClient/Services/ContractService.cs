using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public class ContractService : BaseApiService<LeaseContract>, IContractService {
    public ContractService(IApiClientFactory clients, ISession session, IAppLogger logger)
        : base(clients, session, logger, "/api/contracts") { }
}