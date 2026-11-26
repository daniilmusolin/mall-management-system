using System;
using System.Collections.Generic;
using MallClient.Models;

namespace MallClient.Schemas;

public class SchemaFactory {
    private readonly Dictionary<Type, IEntitySchema> _schemas = new() {
        { typeof(Premise), new PremiseSchema() },
        { typeof(Tenant), new TenantSchema() },
        { typeof(LeaseContract), new LeaseContractSchema() },
        { typeof(Payment), new PaymentSchema() },
        { typeof(ServiceRequest), new ServiceRequestSchema() },
        { typeof(User), new UserSchema() }
    };

    public IEntitySchema GetSchema<T>() where T : class {
        if (_schemas.TryGetValue(typeof(T), out var schema)) return schema;
        throw new NotSupportedException($"Схема для {typeof(T).Name} не найдена");
    }
}