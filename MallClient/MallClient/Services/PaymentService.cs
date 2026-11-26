using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public class PaymentService : BaseApiService<Payment>, IPaymentService {
    public PaymentService(IApiClientFactory clients, ISession session, IAppLogger logger)
        : base(clients, session, logger, "/api/payments") { }
}