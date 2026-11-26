using System.ComponentModel;
using System.Windows;
using System.Windows.Input;
using MallClient.ViewModels;

namespace MallClient.Views;

public partial class LoginWindow : Window {
    public bool IsSwitchingWindow { get; set; } = false; 

    public LoginWindow(LoginViewModel viewModel) {
        InitializeComponent();
        DataContext = viewModel;
    }

    private void LoginWindow_OnClosing(object? sender, CancelEventArgs e) {
        if (!IsSwitchingWindow && Application.Current.MainWindow == this) {
            Application.Current.Shutdown(0);
        }
    }

    private void PasswordBox_OnPasswordChanged(object sender, RoutedEventArgs e) {
        if (DataContext is LoginViewModel vm)
            vm.Password = PasswordBox.Password;
    }

    private void TitleBar_MouseLeftButtonDown(object sender, MouseButtonEventArgs e) {
        if (e.ButtonState == MouseButtonState.Pressed)
            DragMove();
    }

    private void CloseButton_Click(object sender, RoutedEventArgs e) {
        Application.Current.Shutdown(0);
    }
}