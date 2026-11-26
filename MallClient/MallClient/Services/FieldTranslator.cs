using MallClient.Interfaces;

namespace MallClient.Services;

public class FieldTranslator : IFieldTranslator {
    private static readonly Dictionary<string, string> Map = new() {
        ["email"] = "Email",
        ["password"] = "Пароль",
        ["full_name"] = "ФИО",
    };

    public string Translate(string fieldName)
        => Map.TryGetValue(fieldName, out var value) ? value : fieldName;
}