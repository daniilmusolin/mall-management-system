using System;
using System.Threading.Tasks;
using System.Windows.Input;
using MallClient.Interfaces;

namespace MallClient.ViewModels;

public class LoginViewModel : BaseViewModel {
    private readonly IAuthService _auth;
    private readonly IWindowService _windows;

    private string _email = string.Empty;
    private string _password = string.Empty;
    private string _error = string.Empty;

    public string Email { get => _email; set { _email = value; OnPropertyChanged(); } }
    public string Password { get => _password; set { _password = value; OnPropertyChanged(); } }
    public string Error { get => _error; set { _error = value; OnPropertyChanged(); } }

    public ICommand LoginCommand { get; }
    public ICommand GoToRegisterCommand { get; }

    public LoginViewModel(IAuthService auth, IWindowService windows) {
        _auth = auth;
        _windows = windows;
        LoginCommand = new RelayCommand(async _ => await LoginAsync());
        GoToRegisterCommand = new RelayCommand(_ => _windows.ShowRegisterWindow());
    }

    private async Task LoginAsync() {
        try {
            if (string.IsNullOrWhiteSpace(Email) || string.IsNullOrWhiteSpace(Password)) {
                Error = "Заполните все поля";
                return;
            }

            var success = await _auth.LoginAsync(Email, Password);
            if (!success) {
                Error = "Неверный email или пароль";
                return;
            }

            _windows.ShowMainWindow();
        } catch (Exception ex) {
            Error = $"Ошибка входа: {ex.Message}";
        }
    }
}