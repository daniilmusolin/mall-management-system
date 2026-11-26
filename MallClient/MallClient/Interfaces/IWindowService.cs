namespace MallClient.Interfaces;

public interface IWindowService {
    void ShowMainWindow();
    void ShowLoginWindow();
    void ShowRegisterWindow();
    bool ShowEditWindow<T>(T? entity, string title) where T : class;
    bool ShowConfirmWindow(string message);
}