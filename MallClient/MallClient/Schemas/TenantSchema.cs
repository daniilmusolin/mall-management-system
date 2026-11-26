using System.Collections.Generic;
using MallClient.Models;

namespace MallClient.Schemas;

public class TenantSchema : BaseSchema {
    public override string GetLabel(string propertyName) => propertyName switch {
        nameof(Tenant.Name) => "Название",
        nameof(Tenant.Inn) => "ИНН",
        nameof(Tenant.LegalAddress) => "Юридический адрес",
        nameof(Tenant.ContactPhone) => "Контактный телефон",
        nameof(Tenant.Email) => "Email",
        _ => propertyName
    };

    public override void Configure(FieldItem field) {
        if (field.Name == nameof(Tenant.Name)) field.MaxLength = 255;
        else if (field.Name == nameof(Tenant.LegalAddress)) field.MaxLength = 500;
        else if (field.Name == nameof(Tenant.ContactPhone)) field.MaxLength = 16;
        else if (field.Name == nameof(Tenant.Email)) field.MaxLength = 254;
        else if (field.Name == nameof(Tenant.Inn)) {
            field.IsNumeric = true;
            field.MaxLength = 12;
        }

        if (field.Name == nameof(Tenant.ContactPhone)) field.IsPhone = true;
        if (field.Name == nameof(Tenant.Email)) field.IsEmail = true;
    }
}