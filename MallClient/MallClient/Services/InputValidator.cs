using System.Text.RegularExpressions;
using MallClient.Interfaces;
using MallClient.Models;

namespace MallClient.Services;

public class InputValidator : IInputValidator {
    private static readonly Regex EmailRegex = new(@"^[^@\s]+@[^@\s]+\.[^@\s]+$");
    private static readonly Regex PhoneRegex = new(@"^\+?[0-9]{7,15}$");
    private static readonly Regex InnRegex = new(@"^\d{10}$|^\d{12}$");
    private static readonly Regex RoomNumRegex = new(@"^[A-Za-zА-Яа-я0-9\-/]{1,20}$");

    public string? ValidateField(FieldItem field, bool isCreate) {
        if (field.MaxLength > 0 && field.Value.Length > field.MaxLength)
            return $"Поле '{field.Label}': максимум {field.MaxLength} символов.";

        if (isCreate && string.IsNullOrWhiteSpace(field.Value) && !field.IsDropdown) {
            return $"Поле '{field.Label}': обязательное поле.";
        }

        if (field.Name == "Email" && !string.IsNullOrWhiteSpace(field.Value) && !EmailRegex.IsMatch(field.Value))
            return $"Поле '{field.Label}': некорректный email.";

        if (field.Name == "ContactPhone" && !string.IsNullOrWhiteSpace(field.Value) && !PhoneRegex.IsMatch(field.Value))
            return $"Поле '{field.Label}': некорректный телефон (от 7 до 15 цифр).";

        if (field.Name == "Inn" && !string.IsNullOrWhiteSpace(field.Value) && !InnRegex.IsMatch(field.Value))
            return $"Поле '{field.Label}': ИНН должен содержать 10 или 12 цифр.";

        if (field.Name == "RoomNumber" && !string.IsNullOrWhiteSpace(field.Value) && !RoomNumRegex.IsMatch(field.Value))
            return $"Поле '{field.Label}': допустимы буквы, цифры, дефис и слэш (1–20 символов).";

        if (!string.IsNullOrWhiteSpace(field.Value)) {
            int len = field.Value.Length;

            if ((field.Name == "Name" || field.Name == "FullName") && len < 2)
                return $"Поле '{field.Label}': минимум 2 символа.";

            if (field.Name == "LegalAddress" && len < 5)
                return $"Поле '{field.Label}': минимум 5 символов.";

            if (field.Name == "Description" && len < 5)
                return $"Поле '{field.Label}': минимум 5 символов.";

            if (field.Name == "RequestType" && len < 2)
                return $"Поле '{field.Label}': минимум 2 символа.";

            if (field.Name == "Type" && len < 1)
                return $"Поле '{field.Label}': минимум 1 символ.";
        }

        if (field.IsNumeric && !string.IsNullOrWhiteSpace(field.Value)) {
            if (double.TryParse(field.Value.Replace(',', '.'), out var num)) {
                if (field.Name == "Floor") {
                    if (num < 0) return "Этаж не может быть отрицательным.";
                    if (num > 100) return "Этаж не может быть больше 100.";
                }

                if (field.Name == "Area") {
                    if (num <= 0) return "Площадь должна быть больше нуля.";
                    if (num > 9999999) return "Площадь не может превышать 9 999 999.";
                }

                if (field.Name == "Amount") {
                    if (num <= 0) return "Сумма должна быть больше нуля.";
                    if (num > 999999999) return "Сумма не может превышать 999 999 999.";
                }

                if (field.Name == "RatePerSqm") {
                    if (num <= 0) return "Ставка должна быть больше нуля.";
                    if (num > 999999) return "Ставка не может превышать 999 999.";
                }

                if (field.Name == "DepositAmount") {
                    if (num < 0) return "Сумма залога не может быть отрицательной.";
                    if (num > 999999999) return "Залог не может превышать 999 999 999.";
                }
            } else {
                return $"Поле '{field.Label}': введите корректное число.";
            }
        }

        return null;
    }
}