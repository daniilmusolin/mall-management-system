using System;
using System.Collections.ObjectModel;
using System.Threading.Tasks;
using MallClient.Interfaces;
using MallClient.Services;

namespace MallClient.ViewModels;

public class EntityViewModel<T> where T : class {
    private readonly IApiCrudService<T> _service;
    private readonly IWindowService _windows;
    private readonly string _entityName;
    private readonly Func<T, long> _idSelector;
    private readonly Action<string> _onError;
    private readonly IAppLogger _logger;

    public ObservableCollection<T> Items { get; } = new();

    public EntityViewModel(
        IApiCrudService<T> service,
        IWindowService windows,
        IAppLogger logger,
        string entityName,
        Func<T, long> idSelector,
        Action<string> onError) {
        _service = service;
        _windows = windows;
        _logger = logger;
        _entityName = entityName;
        _idSelector = idSelector;
        _onError = onError;
    }

    public async Task LoadAsync() {
        try {
            var data = await _service.GetAllAsync();
            Items.Clear();
            foreach (var item in data) Items.Add(item);
        } catch (Exception ex) {
            _logger.LogError($"Ошибка загрузки {_entityName}", ex);
            _onError($"Не удалось загрузить {_entityName}: {ex.Message}");
        }
    }

    public void ShowAdd() {
        try {
            if (_windows.ShowEditWindow<T>(null, $"Добавить {_entityName}")) _ = LoadAsync();
        } catch (Exception ex) {
            _logger.LogError($"Ошибка открытия формы добавления {_entityName}", ex);
        }
    }

    public void ShowEdit(T item) {
        try {
            if (_windows.ShowEditWindow(item, $"Редактировать {_entityName}")) _ = LoadAsync();
        } catch (Exception ex) {
            _logger.LogError($"Ошибка открытия формы редактирования {_entityName}", ex);
        }
    }

    public async Task DeleteAsync(T item) {
        var id = _idSelector(item);
        if (!_windows.ShowConfirmWindow($"Вы уверены, что хотите удалить {_entityName}?")) return;

        try {
            await _service.DeleteAsync(id);
            _logger.LogInfo($"Удален {_entityName} с ID {id}");
            await LoadAsync();
        } catch (Exception ex) {
            _logger.LogError($"Ошибка удаления {_entityName} с ID {id}", ex);
            _onError($"Не удалось удалить {_entityName}: {ex.Message}");
        }
    }
}