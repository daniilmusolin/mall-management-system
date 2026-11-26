using System.Collections.Generic;
using MallClient.Models;

namespace MallClient.Schemas;

public class LeaseContractSchema : BaseSchema {
    public override string GetLabel(string propertyName) => propertyName switch {
        nameof(LeaseContract.PremiseId) => "Помещение",
        nameof(LeaseContract.TenantId) => "Арендатор",
        nameof(LeaseContract.StartDate) => "Дата начала",
        nameof(LeaseContract.EndDate) => "Дата окончания",
        nameof(LeaseContract.RatePerSqm) => "Ставка за м²",
        nameof(LeaseContract.DepositAmount) => "Сумма залога",
        nameof(LeaseContract.Status) => "Статус",
        _ => propertyName
    };

    public override void ConfigureDropdowns(FieldItem field) {
        if (field.Name == nameof(LeaseContract.Status)) {
            field.IsDropdown = true;
            field.Options = new List<OptionItem> {
                new() { Value = "active", DisplayName = "Активен" },
                new() { Value = "terminated", DisplayName = "Расторгнут" },
                new() { Value = "expired", DisplayName = "Истек" }
            };
        }
    }
}