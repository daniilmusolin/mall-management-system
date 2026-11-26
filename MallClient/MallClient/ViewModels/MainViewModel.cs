using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using System.Windows.Input;
using



using MallClient.ViewModels;
using MallClient.Interfaces;
using MallClient.

namespace ViewModels;

public class MainViewModel : BaseViewModel {
    private readonly ISession _session;
    private readonly IAuthService _auth;
    private readonly IWindowService _windows;

    public string UserName => _session.CurrentUser.FullName;
    public string UserRole => _session.CurrentUser.Role.ToUpper();

    public bool CanViewPremises => _session.IsAdmin || _session.IsManager || _session.IsAccountant;
    public bool CanViewTenants => _session.IsAdmin || _session.IsManager;
    public bool CanViewContracts => _session.IsAdmin || _session.IsManager || _session.IsAccountant;
    public bool CanViewPayments => _session.IsAdmin || _session.IsAccountant;
    public bool CanViewRequests => _session.IsAdmin || _session.IsTechnician || _session.IsSupport;
    public bool CanViewUsers => _session.IsAdmin;

    public bool CanCreatePremise => _session.IsAdmin;
    public bool CanCreateTenant => _session.IsAdmin || _session.IsManager;
    public bool CanCreateContract => _session.IsAdmin || _session.IsManager;
    public bool CanCreatePayment => _session.IsAdmin || _session.IsAccountant;
    public bool CanCreateRequest => _session.IsAdmin || _session.IsSupport;

    public bool CanDeletePremise => _session.IsAdmin;
    public bool CanDeleteTenant => _session.IsAdmin;
    public bool CanDeleteContract => _session.IsAdmin;
    public bool CanDeletePayment => _session.IsAdmin;
    public bool CanDeleteRequest => _session.IsAdmin;
    public bool CanDeleteUser => _session.IsAdmin;

    public EntityViewModel<Premise> Premises { get; }
    public EntityViewModel<Tenant> Tenants { get; }
    public EntityViewModel<LeaseContract> Contracts { get; }
    public EntityViewModel<Payment> Payments { get; }
    public EntityViewModel<ServiceRequest> Requests { get; }
    public EntityViewModel<User> Users { get; }

    private string _errorMessage = string.Empty;
    public string ErrorMessage {
        get => _errorMessage;
        set { _errorMessage = value; OnPropertyChanged(); }
    }

    public ICommand LogoutCommand { get; }

    public MainViewModel(
        ISession session,
        IAuthService auth,
        IWindowService windows,
        IAppLogger logger,
        IApiCrudService<Premise> premisesService,
        IApiCrudService<Tenant> tenantsService,
        IApiCrudService<LeaseContract> contractsService,
        IApiCrudService<Payment> paymentsService,
        IApiCrudService<ServiceRequest> requestsService,
        IApiCrudService<User> usersService) {

        _session = session;
        _auth = auth;
        _windows = windows;

        var errorHandler = new Action<string>(msg => ErrorMessage = msg);

        Premises = new EntityViewModel<Premise>(premisesService, windows, logger, "помещение", p => p.Id, errorHandler);
        Tenants = new EntityViewModel<Tenant>(tenantsService, windows, logger, "арендатора", t => t.Id, errorHandler);
        Contracts = new EntityViewModel<LeaseContract>(contractsService, windows, logger, "договор", c => c.Id, errorHandler);
        Payments = new EntityViewModel<Payment>(paymentsService, windows, logger, "платеж", p => p.Id, errorHandler);
        Requests = new EntityViewModel<ServiceRequest>(requestsService, windows, logger, "заявку", r => r.Id, errorHandler);
        Users = new EntityViewModel<User>(usersService, windows, logger, "пользователя", u => u.Id, errorHandler);

        LogoutCommand = new RelayCommand(_ => Logout());
    }

    private void Logout() {
        try { _auth.Logout(); } catch { }
        _windows.ShowLoginWindow();
    }

    public async Task LoadDataAsync() {
        ErrorMessage = string.Empty;
        var tasks = new List<Task>();

        if (CanViewPremises) tasks.Add(Premises.LoadAsync());
        if (CanViewTenants) tasks.Add(Tenants.LoadAsync());
        if (CanViewContracts) tasks.Add(Contracts.LoadAsync());
        if (CanViewPayments) tasks.Add(Payments.LoadAsync());
        if (CanViewRequests) tasks.Add(Requests.LoadAsync());
        if (CanViewUsers) tasks.Add(Users.LoadAsync());

        await Task.WhenAll(tasks);
    }
}