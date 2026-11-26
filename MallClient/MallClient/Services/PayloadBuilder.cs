using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Reflection;
using MallClient.Models;

namespace MallClient.Services;

public class PayloadBuilder {
    public Dictionary<string, object> Build<T>(T? entity, IEnumerable<FieldItem> fields) where T : class {
        var payload = new Dictionary<string, object>();

        if (entity != null) {
            foreach (var prop in typeof(T).GetProperties()) {
                if (prop.Name == "Id" || prop.Name == "CreatedAt" || prop.Name == "PasswordHash") continue;

                var val = prop.GetValue(entity);
                if (val != null) {
                    var propName = ToSnakeCase(prop.Name);
                    payload[propName] = val;
                }
            }
        }

        foreach (var field in fields) {
            var prop = typeof(T).GetProperty(field.Name);
            if (prop == null) continue;

            var propName = ToSnakeCase(field.Name);
            var t = prop.PropertyType;
            var raw = field.Value?.Trim() ?? string.Empty;

            if (t == typeof(int) || t == typeof(int?)) {
                if (int.TryParse(raw, NumberStyles.Integer, CultureInfo.InvariantCulture, out var i))
                    payload[propName] = i;
            } else if (t == typeof(long) || t == typeof(long?)) {
                if (long.TryParse(raw, NumberStyles.Integer, CultureInfo.InvariantCulture, out var l))
                    payload[propName] = l;
            } else if (t == typeof(decimal) || t == typeof(decimal?)) {
                if (decimal.TryParse(raw.Replace(',', '.'), NumberStyles.Number, CultureInfo.InvariantCulture, out var m))
                    payload[propName] = m;
            } else if (t == typeof(double) || t == typeof(double?)) {
                if (double.TryParse(raw.Replace(',', '.'), NumberStyles.Float, CultureInfo.InvariantCulture, out var d))
                    payload[propName] = d;
            } else {
                payload[propName] = field.Value ?? string.Empty;
            }
        }

        return payload;
    }

    private string ToSnakeCase(string input) {
        if (string.IsNullOrEmpty(input)) return input;
        return string.Concat(input.Select((x, i) =>
            i > 0 && char.IsUpper(x) ? "_" + x.ToString() : x.ToString())).ToLower();
    }
}