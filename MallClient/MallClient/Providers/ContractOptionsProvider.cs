using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Providers;

public class ContractOptionsProvider : IOptionsProvider {
    private readonly IContractService _service;

    public ContractOptionsProvider(IContractService service) {
        _service = service;
    }

    public async Task<List<OptionItem>> GetOptionsAsync() {
        var contracts = await _service.GetAllAsync();
        return contracts.Select(c => new OptionItem {
            Value = c.Id.ToString(),
            DisplayName = $"Договор №{c.Id} (от {c.StartDate})"
        }).ToList();
    }
}