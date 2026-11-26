using System.ComponentModel;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Input;
using MallClient.Models;
using MallClient.ViewModels;

namespace MallClient.Views;

public partial class MainWindow : Window {
    private MainViewModel? _viewModel;
    private bool _initialized;

    public MainWindow(MainViewModel viewModel) {
        _viewModel = viewModel;
        InitializeComponent();
        DataContext = _viewModel;
    }

    private async void MainWindow_OnLoaded(object sender, RoutedEventArgs e) {
        if (_viewModel == null) return;

        await _viewModel.LoadDataAsync();

        _initialized = true;

        if (TabPremises.IsVisible) TabPremises.IsChecked = true;
        else if (TabTenants.IsVisible) TabTenants.IsChecked = true;
        else if (TabContracts.IsVisible) TabContracts.IsChecked = true;
        else if (TabPayments.IsVisible) TabPayments.IsChecked = true;
        else if (TabRequests.IsVisible) TabRequests.IsChecked = true;
        else if (TabUsers.IsVisible) TabUsers.IsChecked = true;
    }

    private void MainWindow_OnClosing(object? sender, CancelEventArgs e) {
        if (Application.Current.MainWindow == this) {
            App.ForceExit();
        }
    }

    private void TitleBar_MouseLeftButtonDown(object sender, MouseButtonEventArgs e) {
        if (e.ButtonState == MouseButtonState.Pressed)
            DragMove();
    }

    private void CloseButton_Click(object sender, RoutedEventArgs e) {
        App.ForceExit();
    }

    private async void Tab_Checked(object sender, RoutedEventArgs e) {
        if (_viewModel == null) return;
        if (sender is not RadioButton rb) return;

        PanelPremises.Visibility = Visibility.Collapsed;
        PanelTenants.Visibility = Visibility.Collapsed;
        PanelContracts.Visibility = Visibility.Collapsed;
        PanelPayments.Visibility = Visibility.Collapsed;
        PanelRequests.Visibility = Visibility.Collapsed;
        PanelUsers.Visibility = Visibility.Collapsed;

        switch (rb.Name) {
            case nameof(TabPremises):
                PanelPremises.Visibility = Visibility.Visible;
                if (_initialized) await _viewModel.Premises.LoadAsync();
                break;
            case nameof(TabTenants):
                PanelTenants.Visibility = Visibility.Visible;
                if (_initialized) await _viewModel.Tenants.LoadAsync();
                break;
            case nameof(TabContracts):
                PanelContracts.Visibility = Visibility.Visible;
                if (_initialized) await _viewModel.Contracts.LoadAsync();
                break;
            case nameof(TabPayments):
                PanelPayments.Visibility = Visibility.Visible;
                if (_initialized) await _viewModel.Payments.LoadAsync();
                break;
            case nameof(TabRequests):
                PanelRequests.Visibility = Visibility.Visible;
                if (_initialized) await _viewModel.Requests.LoadAsync();
                break;
            case nameof(TabUsers):
                PanelUsers.Visibility = Visibility.Visible;
                if (_initialized) await _viewModel.Users.LoadAsync();
                break;
        }
    }

    private void AddPremise_Click(object sender, RoutedEventArgs e) => _viewModel?.Premises.ShowAdd();
    private void PremisesGrid_MouseDoubleClick(object sender, MouseButtonEventArgs e) {
        if (PremisesGrid.SelectedItem is Premise premise) _viewModel?.Premises.ShowEdit(premise);
    }
    private async void DeletePremise_Click(object sender, RoutedEventArgs e) {
        if (_viewModel == null) return;
        if (sender is Button btn && btn.Tag is long id) {
            var item = _viewModel.Premises.Items.FirstOrDefault(x => x.Id == id);
            if (item != null) await _viewModel.Premises.DeleteAsync(item);
        }
    }

    private void AddTenant_Click(object sender, RoutedEventArgs e) => _viewModel?.Tenants.ShowAdd();
    private void TenantsGrid_MouseDoubleClick(object sender, MouseButtonEventArgs e) {
        if (TenantsGrid.SelectedItem is Tenant tenant) _viewModel?.Tenants.ShowEdit(tenant);
    }
    private async void DeleteTenant_Click(object sender, RoutedEventArgs e) {
        if (_viewModel == null) return;
        if (sender is Button btn && btn.Tag is long id) {
            var item = _viewModel.Tenants.Items.FirstOrDefault(x => x.Id == id);
            if (item != null) await _viewModel.Tenants.DeleteAsync(item);
        }
    }

    private void AddContract_Click(object sender, RoutedEventArgs e) => _viewModel?.Contracts.ShowAdd();
    private void ContractsGrid_MouseDoubleClick(object sender, MouseButtonEventArgs e) {
        if (ContractsGrid.SelectedItem is LeaseContract contract) _viewModel?.Contracts.ShowEdit(contract);
    }
    private async void DeleteContract_Click(object sender, RoutedEventArgs e) {
        if (_viewModel == null) return;
        if (sender is Button btn && btn.Tag is long id) {
            var item = _viewModel.Contracts.Items.FirstOrDefault(x => x.Id == id);
            if (item != null) await _viewModel.Contracts.DeleteAsync(item);
        }
    }

    private void AddPayment_Click(object sender, RoutedEventArgs e) => _viewModel?.Payments.ShowAdd();
    private void PaymentsGrid_MouseDoubleClick(object sender, MouseButtonEventArgs e) {
        if (PaymentsGrid.SelectedItem is Payment payment) _viewModel?.Payments.ShowEdit(payment);
    }
    private async void DeletePayment_Click(object sender, RoutedEventArgs e) {
        if (_viewModel == null) return;
        if (sender is Button btn && btn.Tag is long id) {
            var item = _viewModel.Payments.Items.FirstOrDefault(x => x.Id == id);
            if (item != null) await _viewModel.Payments.DeleteAsync(item);
        }
    }

    private void AddRequest_Click(object sender, RoutedEventArgs e) => _viewModel?.Requests.ShowAdd();
    private void RequestsGrid_MouseDoubleClick(object sender, MouseButtonEventArgs e) {
        if (RequestsGrid.SelectedItem is ServiceRequest request) _viewModel?.Requests.ShowEdit(request);
    }
    private async void DeleteRequest_Click(object sender, RoutedEventArgs e) {
        if (_viewModel == null) return;
        if (sender is Button btn && btn.Tag is long id) {
            var item = _viewModel.Requests.Items.FirstOrDefault(x => x.Id == id);
            if (item != null) await _viewModel.Requests.DeleteAsync(item);
        }
    }

    private void UsersGrid_MouseDoubleClick(object sender, MouseButtonEventArgs e) {
        if (UsersGrid.SelectedItem is User user) _viewModel?.Users.ShowEdit(user);
    }
    private async void DeleteUser_Click(object sender, RoutedEventArgs e) {
        if (_viewModel == null) return;
        if (sender is Button btn && btn.Tag is long id) {
            var item = _viewModel.Users.Items.FirstOrDefault(x => x.Id == id);
            if (item != null) await _viewModel.Users.DeleteAsync(item);
        }
    }
}