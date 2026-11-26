using System;
using System.Windows;
using System.Windows.Threading;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using MallClient.Extensions;
using MallClient.Interfaces;
using MallClient.Models;
using MallClient.Views;

namespace MallClient;

public partial class App : Application {
    private const string AppSettingsFile = "appsettings.json";
    private const string ApiSettingsSection = "ApiSettings";
    private const string ConfigErrorTitle = "Ошибка конфигурации";
    private const string ConfigErrorMessage = "Не найден BaseUrl в appsettings.json";

    public static IServiceProvider Services { get; private set; } = null!;
    private static IAppLogger? _logger;

    protected override void OnStartup(StartupEventArgs e) {
        base.OnStartup(e);

        var configuration = LoadConfiguration();
        var apiSettings = configuration.GetSection(ApiSettingsSection).Get<AppSettings>();

        if (apiSettings == null || string.IsNullOrEmpty(apiSettings.BaseUrl)) {
            MessageBox.Show(ConfigErrorMessage, ConfigErrorTitle, MessageBoxButton.OK, MessageBoxImage.Error);
            Shutdown();
            return;
        }

        var services = new ServiceCollection();
        services.AddAppServices(apiSettings.BaseUrl);
        Services = services.BuildServiceProvider();

        _logger = Services.GetRequiredService<IAppLogger>();

        DispatcherUnhandledException += OnDispatcherUnhandledException;
        AppDomain.CurrentDomain.UnhandledException += OnDomainUnhandledException;
        AppDomain.CurrentDomain.ProcessExit += OnProcessExit;
        TaskScheduler.UnobservedTaskException += OnUnobservedTaskException;
        SessionEnding += OnSessionEnding;

        var session = Services.GetRequiredService<ISession>();
        session.SessionExpired += OnSessionExpired;

        var windowService = Services.GetRequiredService<IWindowService>();
        windowService.ShowLoginWindow();
    }

    private void OnSessionExpired() {
        Current?.Dispatcher?.Invoke(() => {
            try {
                var auth = Services.GetRequiredService<IAuthService>();
                var windowService = Services.GetRequiredService<IWindowService>();

                if (Current.MainWindow is LoginWindow) return;

                auth.Logout();
                windowService.ShowLoginWindow();
                _logger?.LogWarning("Сессия истекла, принудительный выход.");

                MessageBox.Show("Ваш аккаунт был удален или срок действия сессии истек.\nПожалуйста, авторизуйтесь заново.",
                    "Доступ запрещен", MessageBoxButton.OK, MessageBoxImage.Warning);
            } catch (Exception ex) {
                _logger?.LogError("Ошибка при обработке истечения сессии", ex);
            }
        });
    }

    private void OnDispatcherUnhandledException(object sender, DispatcherUnhandledExceptionEventArgs e) {
        _logger?.LogError("Необработанная ошибка UI", e.Exception);
        MessageBox.Show(e.Exception.ToString(), "Необработанная ошибка", MessageBoxButton.OK, MessageBoxImage.Error);
        e.Handled = true;
    }

    private void OnDomainUnhandledException(object sender, UnhandledExceptionEventArgs e) {
        var ex = e.ExceptionObject as Exception;
        _logger?.LogError("Критическая ошибка домена", ex);
    }

    private void OnUnobservedTaskException(object? sender, UnobservedTaskExceptionEventArgs e) {
        _logger?.LogError("Необработанная ошибка в Task", e.Exception);
        e.SetObserved();
    }

    private void OnSessionEnding(object sender, SessionEndingCancelEventArgs e) => ForceExit();

    private void OnProcessExit(object? sender, EventArgs e) => DisposeServices();

    protected override void OnExit(ExitEventArgs e) {
        DisposeServices();
        base.OnExit(e);
    }

    public static void ForceExit() {
        try { DisposeServices(); } catch { }
        try { Current?.Shutdown(0); } catch { }
        Environment.Exit(0);
    }

    private static void DisposeServices() {
        if (Services is IDisposable disposable) {
            try { disposable.Dispose(); } catch { }
        }
    }

    private IConfiguration LoadConfiguration() {
        return new ConfigurationBuilder()
            .SetBasePath(AppContext.BaseDirectory)
            .AddJsonFile(AppSettingsFile, optional: false, reloadOnChange: true)
            .Build();
    }
}