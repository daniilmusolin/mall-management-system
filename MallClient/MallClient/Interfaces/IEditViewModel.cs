using MallClient.Models;
using System.Collections.ObjectModel;
using System.Windows.Input;

namespace MallClient.Interfaces;

public interface IEditViewModel {
    ObservableCollection<FieldItem> Fields { get; }
    string ErrorMessage { get; }
    ICommand SaveCommand { get; }
}