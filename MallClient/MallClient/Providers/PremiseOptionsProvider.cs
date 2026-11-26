using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Providers;

public class PremiseOptionsProvider : IOptionsProvider {
    private readonly IPremiseService _service;

    public PremiseOptionsProvider(IPremiseService service) {
        _service = service;
    }

    public async Task<List<OptionItem>> GetOptionsAsync() {
        var premises = await _service.GetAllAsync();
        return premises.Select(p => new OptionItem {
            Value = p.Id.ToString(),
            DisplayName = $"ID {p.Id} - {p.RoomNumber} (этаж {p.Floor})"
        }).ToList();
    }
}