using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public class UserService : BaseApiService<User>, IUserService {
    public UserService(IApiClientFactory clients, ISession session, IAppLogger logger)
        : base(clients, session, logger, "/api/users") { }
}