using MallClient.Models;

namespace MallClient.Interfaces;

public interface IInputValidator {
    string? ValidateField(FieldItem field, bool isCreate);
}