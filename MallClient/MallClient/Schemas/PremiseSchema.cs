using System.Collections.Generic;
using MallClient.Models;

namespace MallClient.Schemas;

public class PremiseSchema : BaseSchema {
    public override string GetLabel(string propertyName) => propertyName switch {
        nameof(Premise.RoomNumber) => "Номер помещения",
        nameof(Premise.Floor) => "Этаж",
        nameof(Premise.Area) => "Площадь",
        nameof(Premise.Type) => "Тип",
        nameof(Premise.Zone) => "Зона",
        nameof(Premise.Status) => "Статус",
        _ => propertyName
    };

    public override void Configure(FieldItem field) {
        if (field.Name == nameof(Premise.RoomNumber)) field.MaxLength = 20;
        else if (field.Name == nameof(Premise.Type)) field.MaxLength = 50;
        else if (field.Name == nameof(Premise.Zone)) field.MaxLength = 20;
    }

    public override void ConfigureDropdowns(FieldItem field) {
        if (field.Name == nameof(Premise.Status)) {
            field.IsDropdown = true;
            field.Options = new List<OptionItem> {
                new() { Value = "free", DisplayName = "Свободно" },
                new() { Value = "rented", DisplayName = "Арендовано" },
                new() { Value = "maintenance", DisplayName = "На ремонте" }
            };
        }
    }
}