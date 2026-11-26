using System.Threading.Tasks;

namespace MallClient.Interfaces;

public interface IAuthService {
    Task<bool> LoginAsync(string email, string password);
    Task<(bool Success, string Error)> RegisterAsync(string email, string password, string fullName);
    void Logout();
}