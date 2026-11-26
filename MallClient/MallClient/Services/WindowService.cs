using System;
using System.Windows;
using Microsoft.Extensions.DependencyInjection;
using MallClient.Interfaces;
using MallClient.Models;
using MallClient.ViewModels;
using MallClient.Views;
using MallClient.Schemas;

namespace MallClient.Services;

public class WindowService : IWindowService {
    private readonly IServiceProvider _serviceProvider;
    private readonly SchemaFactory _schemaFactory;

    public WindowService(IServiceProvider serviceProvider, SchemaFactory schemaFactory) {
        _serviceProvider = serviceProvider;
        _schemaFactory = schemaFactory;
    }

    public void ShowLoginWindow() => SwitchWindow(_serviceProvider.GetRequiredService<LoginWindow>());
    public void ShowRegisterWindow() => SwitchWindow(_serviceProvider.GetRequiredService<RegisterWindow>());
    public void ShowMainWindow() => SwitchWindow(_serviceProvider.GetRequiredService<MainWindow>());

    private void SwitchWindow(Window newWindow) {
        var oldWindow = Application.Current.MainWindow;

        Application.Current.MainWindow = newWindow;
        newWindow.Show();

        if (oldWindow == null || oldWindow == newWindow) return;

        if (oldWindow is LoginWindow lw) lw.IsSwitchingWindow = true;
        if (oldWindow is RegisterWindow rw) rw.IsSwitchingWindow = true;

        if (oldWindow.IsLoaded) oldWindow.Close();
    }

    public bool ShowEditWindow<T>(T? entity, string title) where T : class {
        var service = GetCrudService<T>();
        var validator = _serviceProvider.GetRequiredService<IInputValidator>();
        var payloadBuilder = _serviceProvider.GetRequiredService<PayloadBuilder>();
        var schema = _schemaFactory.GetSchema<T>();

        var vm = new GenericEditViewModel<T>(
            entity,
            service,
            validator,
            payloadBuilder,
            schema,
            _serviceProvider
        );

        var win = _serviceProvider.GetRequiredService<GenericEditWindow>();
        win.DataContext = vm;
        win.Title = title;
        win.Owner = Application.Current.MainWindow;

        vm.CloseAction = () => win.DialogResult = true;

        return win.ShowDialog() == true;
    }

    public bool ShowConfirmWindow(string message) {
        var win = new ConfirmWindow(message) {
            Owner = Application.Current.MainWindow
        };
        return win.ShowDialog() == true;
    }

    private IApiCrudService<T> GetCrudService<T>() where T : class {
        return typeof(T).Name switch {
            nameof(Premise) => (IApiCrudService<T>)(object)_serviceProvider.GetRequiredService<IPremiseService>(),
            nameof(Tenant) => (IApiCrudService<T>)(object)_serviceProvider.GetRequiredService<ITenantService>(),
            nameof(LeaseContract) => (IApiCrudService<T>)(object)_serviceProvider.GetRequiredService<IContractService>(),
            nameof(Payment) => (IApiCrudService<T>)(object)_serviceProvider.GetRequiredService<IPaymentService>(),
            nameof(ServiceRequest) => (IApiCrudService<T>)(object)_serviceProvider.GetRequiredService<IRequestService>(),
            nameof(User) => (IApiCrudService<T>)(object)_serviceProvider.GetRequiredService<IUserService>(),
            _ => throw new NotSupportedException($"Для типа {typeof(T).Name} не настроен IApiCrudService")
        };
    }
}