using System.Collections.Generic;
using MallClient.Models;

namespace MallClient.Schemas;

public class PaymentSchema : BaseSchema {
    public override string GetLabel(string propertyName) => propertyName switch {
        nameof(Payment.ContractId) => "Договор",
        nameof(Payment.Amount) => "Сумма",
        nameof(Payment.PaymentDate) => "Дата платежа",
        nameof(Payment.PaymentType) => "Тип платежа",
        nameof(Payment.Status) => "Статус",
        nameof(Payment.ReferenceNumber) => "Номер документа",
        _ => propertyName
    };

    public override void Configure(FieldItem field) {
        if (field.Name == nameof(Payment.ReferenceNumber)) field.MaxLength = 100;
    }

    public override void ConfigureDropdowns(FieldItem field) {
        if (field.Name == nameof(Payment.PaymentType)) {
            field.IsDropdown = true;
            field.Options = new List<OptionItem> {
                new() { Value = "rent", DisplayName = "Аренда" },
                new() { Value = "deposit", DisplayName = "Залог" },
                new() { Value = "utilities", DisplayName = "Коммунальные" },
                new() { Value = "penalty", DisplayName = "Штраф" }
            };
        } else if (field.Name == nameof(Payment.Status)) {
            field.IsDropdown = true;
            field.Options = new List<OptionItem> {
                new() { Value = "pending", DisplayName = "Ожидает" },
                new() { Value = "completed", DisplayName = "Завершен" },
                new() { Value = "cancelled", DisplayName = "Отменен" }
            };
        }
    }
}