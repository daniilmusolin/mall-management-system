using MallClient.Models;

namespace MallClient.Schemas;

public interface IEntitySchema {
    string GetLabel(string propertyName);
    void Configure(FieldItem field);
    void ConfigureDropdowns(FieldItem field);
    List<OptionItem>? GetOptions(string propertyName);
}