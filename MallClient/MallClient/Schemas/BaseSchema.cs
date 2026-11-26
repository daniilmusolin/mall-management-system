using System;
using System.Collections.Generic;
using MallClient.Models;

namespace MallClient.Schemas;

public abstract class BaseSchema : IEntitySchema {
    public virtual string GetLabel(string propertyName) => propertyName;

    public virtual void Configure(FieldItem field) {
        // По умолчанию ничего не делаем
    }

    public virtual void ConfigureDropdowns(FieldItem field) {
        // По умолчанию ничего не делаем
    }

    public virtual List<OptionItem>? GetOptions(string propertyName) => null;
}