using MallClient.Interfaces;
using MallClient.Models;
using MallClient.Providers;
using MallClient.Schemas;
using MallClient.Services;
using MallClient.ViewModels;
using MallClient.Views;
using Microsoft.Extensions.DependencyInjection;

namespace MallClient.Extensions;

public static class ServiceCollectionExtensions {
    public static IServiceCollection AddAppServices(this IServiceCollection services, string baseUrl) {
        services.AddHttpClient("MallClient", c => {
            c.BaseAddress = new Uri(baseUrl);
            c.Timeout = TimeSpan.FromSeconds(30);
        })
        .AddHttpMessageHandler<AuthMessageHandler>()
        .AddHttpMessageHandler<LoggingMessageHandler>();

        services.AddSingleton<IAppLogger, FileLogger>();
        services.AddTransient<LoggingMessageHandler>();
        services.AddTransient<AuthMessageHandler>();

        services.AddSingleton<ISession, Session>();
        services.AddSingleton<ITokenStore, TokenStore>();
        services.AddSingleton<IApiClientFactory, ApiClientFactory>();
        services.AddSingleton<IFieldTranslator, FieldTranslator>();
        services.AddSingleton<IApiErrorParser, ApiErrorParser>();
        services.AddSingleton<IAuthApi, AuthApi>();
        services.AddSingleton<IAuthService, AuthService>();

        services.AddSingleton<IApiCrudService<Premise>, PremiseService>();
        services.AddSingleton<IApiCrudService<Tenant>, TenantService>();
        services.AddSingleton<IApiCrudService<LeaseContract>, ContractService>();
        services.AddSingleton<IApiCrudService<Payment>, PaymentService>();
        services.AddSingleton<IApiCrudService<ServiceRequest>, RequestService>();
        services.AddSingleton<IApiCrudService<User>, UserService>();

        services.AddSingleton<IWindowService, WindowService>();
        services.AddSingleton<IInputValidator, InputValidator>();
        services.AddSingleton<PayloadBuilder>();
        services.AddSingleton<SchemaFactory>();

        services.AddTransient<PremiseOptionsProvider>();
        services.AddTransient<TenantOptionsProvider>();
        services.AddTransient<ContractOptionsProvider>();

        services.AddTransient<LoginViewModel>();
        services.AddTransient<RegisterViewModel>();
        services.AddTransient<MainViewModel>();
        services.AddTransient<LoginWindow>();
        services.AddTransient<RegisterWindow>();
        services.AddTransient<MainWindow>();
        services.AddTransient<GenericEditWindow>();

        return services;
    }
}