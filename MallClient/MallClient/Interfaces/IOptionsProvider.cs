using System.Collections.Generic;
using System.Threading.Tasks;
using MallClient.Models;

namespace MallClient.Interfaces;

public interface IOptionsProvider {
    Task<List<OptionItem>> GetOptionsAsync();
}