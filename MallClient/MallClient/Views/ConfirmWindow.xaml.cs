using System.Windows;
using System.Windows.Input;

namespace MallClient.Views;

public partial class ConfirmWindow : Window {
    public ConfirmWindow(string message) {
        InitializeComponent();
        MessageText.Text = message;
    }

    private void TitleBar_MouseLeftButtonDown(object sender, MouseButtonEventArgs e) {
        if (e.ButtonState == MouseButtonState.Pressed)
            DragMove();
    }

    private void Confirm_Click(object sender, RoutedEventArgs e) {
        DialogResult = true;
    }

    private void Cancel_Click(object sender, RoutedEventArgs e) {
        DialogResult = false;
    }
}