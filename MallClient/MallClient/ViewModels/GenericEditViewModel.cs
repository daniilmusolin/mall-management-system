using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.Linq;
using System.Reflection;
using System.Threading.Tasks;
using System.Windows.Input;
using Microsoft.Extensions.DependencyInjection;
using MallClient.Interfaces;
using MallClient.Models;
using MallClient.Schemas;
using MallClient.Services;

namespace MallClient.ViewModels;

public class GenericEditViewModel<T> : BaseViewModel, IEditViewModel where T : class {
    private readonly IApiCrudService<T> _service;
    private readonly IInputValidator _validator;
    private readonly PayloadBuilder _payloadBuilder;
    private readonly IEntitySchema _schema;
    private readonly IServiceProvider _serviceProvider;
    private readonly T? _entity;
    private readonly long _entityId;

    private string _errorMessage = string.Empty;
    public string ErrorMessage { get => _errorMessage; set { _errorMessage = value; OnPropertyChanged(); } }

    public ObservableCollection<FieldItem> Fields { get; } = new();

    private bool _canSave = true;
    public bool CanSave { get => _canSave; set { _canSave = value; OnPropertyChanged(); } }

    public ICommand SaveCommand { get; }
    public bool IsSaved { get; private set; }
    public Action? CloseAction { get; set; }

    public GenericEditViewModel(
        T? entity,
        IApiCrudService<T> service,
        IInputValidator validator,
        PayloadBuilder payloadBuilder,
        IEntitySchema schema,
        IServiceProvider serviceProvider) {
        _entity = entity;
        _service = service;
        _validator = validator;
        _payloadBuilder = payloadBuilder;
        _schema = schema;
        _serviceProvider = serviceProvider;

        SaveCommand = new RelayCommand(async _ => await SaveAsync(), _ => CanSave);

        if (_entity != null) {
            var idProp = typeof(T).GetProperty("Id");
            _entityId = (long)(idProp?.GetValue(_entity) ?? 0);
        }

        LoadFields();
    }

    private void LoadFields() {
        foreach (var prop in typeof(T).GetProperties()) {
            if (prop.Name == "Id" || prop.Name == "PasswordHash" || prop.Name == "CreatedAt") continue;

            var val = _entity != null ? prop.GetValue(_entity)?.ToString() : string.Empty;
            var field = new FieldItem { Name = prop.Name, Label = _schema.GetLabel(prop.Name), Value = val ?? string.Empty };

            var t = prop.PropertyType;
            var fkAttr = prop.GetCustomAttribute<ForeignKeyAttribute>();

            if (fkAttr != null) {
                field.IsDropdown = true;
                field.IsNumeric = false;
            } else if (t == typeof(int) || t == typeof(long) || t == typeof(double) || t == typeof(decimal)) {
                field.IsNumeric = true;
                if (field.MaxLength == 0) field.MaxLength = 10;
            }

            if (prop.Name.EndsWith("Date", StringComparison.OrdinalIgnoreCase)) {
                field.IsDate = true;
                field.IsNumeric = false;
            }

            _schema.Configure(field);
            _schema.ConfigureDropdowns(field);

            Fields.Add(field);

            if (fkAttr != null) {
                _ = LoadDropdownOptionsAsync(field, fkAttr.ProviderType);
            }
        }
    }

    private async Task LoadDropdownOptionsAsync(FieldItem field, Type providerType) {
        field.Options.Add(new OptionItem { Value = string.Empty, DisplayName = "Загрузка..." });

        try {
            var provider = (IOptionsProvider)_serviceProvider.GetRequiredService(providerType);
            var options = await provider.GetOptionsAsync();

            field.Options.Clear();
            foreach (var opt in options) {
                field.Options.Add(opt);
            }

            if (string.IsNullOrEmpty(field.Value) && field.Options.Any()) {
                field.Value = field.Options.First().Value;
            }
        } catch (Exception ex) {
            field.Options.Clear();
            field.Options.Add(new OptionItem { Value = "", DisplayName = $"Ошибка загрузки: {ex.Message}" });
            System.Diagnostics.Debug.WriteLine($"[FK ERROR] {ex}");
        }
    }

    private async Task SaveAsync() {
        try {
            ErrorMessage = string.Empty;
            bool isCreate = _entity == null;

            foreach (var field in Fields) {
                var error = _validator.ValidateField(field, isCreate);
                if (error != null) {
                    ErrorMessage = error;
                    return;
                }
            }

            var payload = _payloadBuilder.Build(_entity, Fields);
            if (isCreate) await _service.CreateAsync(payload);
            else await _service.UpdateAsync(_entityId, payload);

            IsSaved = true;
            CloseAction?.Invoke();
        } catch (Exception ex) {
            ErrorMessage = ex.Message;
        }
    }
}