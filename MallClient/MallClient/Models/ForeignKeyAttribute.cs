using System;

namespace MallClient.Models;

[AttributeUsage(AttributeTargets.Property)]
public class ForeignKeyAttribute : Attribute {
    public Type ProviderType { get; }

    public ForeignKeyAttribute(Type providerType) {
        ProviderType = providerType;
    }
}