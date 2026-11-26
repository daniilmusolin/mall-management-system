using System.Text.RegularExpressions;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using MallClient.Models;

namespace MallClient.Views;

public partial class GenericEditWindow : Window {
    private static readonly Regex NumericRegex = new(@"^[0-9]*[.,]?[0-9]*$");
    private static readonly Regex PhoneRegex = new(@"^\+?[0-9]*$");

    public GenericEditWindow() {
        InitializeComponent();
    }

    private void TitleBar_MouseLeftButtonDown(object sender, MouseButtonEventArgs e) {
        if (e.ButtonState == MouseButtonState.Pressed)
            DragMove();
    }

    private void CloseButton_Click(object sender, RoutedEventArgs e) {
        DialogResult = false;
    }

    private void NumericTextBox_PreviewTextInput(object sender, TextCompositionEventArgs e) {
        if (sender is not TextBox tb) return;
        if (tb.DataContext is not FieldItem field) return;

        if (field.MaxLength > 0 && tb.Text.Length >= field.MaxLength) {
            e.Handled = true;
            return;
        }

        var proposed = tb.Text.Insert(tb.CaretIndex, e.Text);
        if (!NumericRegex.IsMatch(proposed))
            e.Handled = true;
    }

    private void NumericTextBox_Pasting(object sender, DataObjectPastingEventArgs e) {
        if (e.DataObject.GetDataPresent(typeof(string))) {
            var text = (string)e.DataObject.GetData(typeof(string));
            if (!NumericRegex.IsMatch(text))
                e.CancelCommand();
        } else {
            e.CancelCommand();
        }
    }

    private void PhoneTextBox_PreviewTextInput(object sender, TextCompositionEventArgs e) {
        if (sender is not TextBox tb) return;
        if (tb.DataContext is not FieldItem field) return;

        if (field.MaxLength > 0 && tb.Text.Length >= field.MaxLength) {
            e.Handled = true;
            return;
        }

        var proposed = tb.Text.Insert(tb.CaretIndex, e.Text);
        if (!PhoneRegex.IsMatch(proposed))
            e.Handled = true;
    }
}