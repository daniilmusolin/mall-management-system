using System.Collections.Generic;
using MallClient.Models;

namespace MallClient.Schemas;

public class UserSchema : BaseSchema {
    public override string GetLabel(string propertyName) => propertyName switch {
        nameof(User.Email) => "Email",
        nameof(User.FullName) => "ФИО",
        nameof(User.Role) => "Роль",
        _ => propertyName
    };

    public override void Configure(FieldItem field) {
        if (field.Name == nameof(User.Email)) {
            field.MaxLength = 254;
            field.IsEmail = true;
        }
        if (field.Name == nameof(User.FullName)) field.MaxLength = 255;
    }

    public override void ConfigureDropdowns(FieldItem field) {
        if (field.Name == nameof(User.Role)) {
            field.IsDropdown = true;
            field.Options = new List<OptionItem> {
                new() { Value = "admin", DisplayName = "Администратор" },
                new() { Value = "manager", DisplayName = "Менеджер" },
                new() { Value = "accountant", DisplayName = "Бухгалтер" },
                new() { Value = "technician", DisplayName = "Техник" },
                new() { Value = "support", DisplayName = "Поддержка" }
            };
        }
    }
}