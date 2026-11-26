using System.Collections.Generic;
using MallClient.Models;

namespace MallClient.Schemas;

public class ServiceRequestSchema : BaseSchema {
    public override string GetLabel(string propertyName) => propertyName switch {
        nameof(ServiceRequest.PremiseId) => "Помещение",
        nameof(ServiceRequest.RequestType) => "Тип заявки",
        nameof(ServiceRequest.Description) => "Описание",
        nameof(ServiceRequest.Priority) => "Приоритет",
        nameof(ServiceRequest.Status) => "Статус",
        _ => propertyName
    };

    public override void Configure(FieldItem field) {
        if (field.Name == nameof(ServiceRequest.RequestType)) field.MaxLength = 100;
        else if (field.Name == nameof(ServiceRequest.Description)) field.MaxLength = 2000;
    }

    public override void ConfigureDropdowns(FieldItem field) {
        if (field.Name == nameof(ServiceRequest.Priority)) {
            field.IsDropdown = true;
            field.Options = new List<OptionItem> {
                new() { Value = "low", DisplayName = "Низкий" },
                new() { Value = "medium", DisplayName = "Средний" },
                new() { Value = "high", DisplayName = "Высокий" },
                new() { Value = "critical", DisplayName = "Критический" }
            };
        } else if (field.Name == nameof(ServiceRequest.Status)) {
            field.IsDropdown = true;
            field.Options = new List<OptionItem> {
                new() { Value = "open", DisplayName = "Открыта" },
                new() { Value = "in_progress", DisplayName = "В работе" },
                new() { Value = "closed", DisplayName = "Закрыта" },
                new() { Value = "cancelled", DisplayName = "Отменена" }
            };
        }
    }
}