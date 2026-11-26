using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Providers;

public class TenantOptionsProvider : IOptionsProvider {
    private readonly ITenantService _service;

    public TenantOptionsProvider(ITenantService service) {
        _service = service;
    }

    public async Task<List<OptionItem>> GetOptionsAsync() {
        var tenants = await _service.GetAllAsync();
        return tenants.Select(t => new OptionItem {
            Value = t.Id.ToString(),
            DisplayName = $"ID {t.Id} - {t.Name}"
        }).ToList();
    }
}