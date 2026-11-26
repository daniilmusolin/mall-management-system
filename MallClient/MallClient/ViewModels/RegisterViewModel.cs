using System;
using System.Threading.Tasks;
using System.Windows.Input;
using MallClient.Interfaces;

namespace MallClient.ViewModels;

public class RegisterViewModel : BaseViewModel {
    private readonly IAuthService _auth;
    private readonly IWindowService _windows;

    private string _email = string.Empty;
    private string _password = string.Empty;
    private string _fullName = string.Empty;
    private string _error = string.Empty;

    public string Email { get => _email; set { _email = value; OnPropertyChanged(); } }
    public string Password { get => _password; set { _password = value; OnPropertyChanged(); } }
    public string FullName { get => _fullName; set { _fullName = value; OnPropertyChanged(); } }
    public string Error { get => _error; set { _error = value; OnPropertyChanged(); } }

    public ICommand RegisterCommand { get; }
    public ICommand GoToLoginCommand { get; }

    public RegisterViewModel(IAuthService auth, IWindowService windows) {
        _auth = auth;
        _windows = windows;

        RegisterCommand = new RelayCommand(async _ => await RegisterAsync());
        GoToLoginCommand = new RelayCommand(_ => _windows.ShowLoginWindow());
    }

    private async Task RegisterAsync() {
        try {
            Error = string.Empty;

            if (string.IsNullOrWhiteSpace(Email) || string.IsNullOrWhiteSpace(Password) || string.IsNullOrWhiteSpace(FullName)) {
                Error = "Заполните все поля";
                return;
            }

            var (success, error) = await _auth.RegisterAsync(Email, Password, FullName);
            if (!success) {
                Error = error;
                return;
            }

            _windows.ShowLoginWindow();
        } catch (Exception ex) {
            Error = $"Ошибка регистрации: {ex.Message}";
        }
    }
}